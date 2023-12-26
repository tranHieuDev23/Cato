package logic

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/gammazero/workerpool"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc/rpcclient"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type Judge interface {
	JudgeSubmission(ctx context.Context, problem *db.Problem, language string, content string) (db.SubmissionResult, error)
	ScheduleJudgeLocalSubmission(id uint64)
	ScheduleJudgeDistributedFirstSubmittedSubmission()
}

type judge struct {
	problemDataAccessor             db.ProblemDataAccessor
	submissionDataAccessor          db.SubmissionDataAccessor
	testCaseDataAccessor            db.TestCaseDataAccessor
	problemTestCaseHashDataAccessor db.ProblemTestCaseHashDataAccessor
	db                              *gorm.DB
	apiClient                       rpcclient.APIClient
	logger                          *zap.Logger
	logicConfig                     configs.Logic
	appArguments                    utils.Arguments
	workerPool                      *workerpool.WorkerPool
	languageToCompileLogic          map[string]Compile
	languageToTestCaseRunLogic      map[string]TestCaseRun
	submissionRetryDelayDuration    time.Duration
}

func NewJudge(
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	testCaseDataAccessor db.TestCaseDataAccessor,
	problemTestCaseHashDataAccessor db.ProblemTestCaseHashDataAccessor,
	dockerClient *client.Client,
	db *gorm.DB,
	apiClient rpcclient.APIClient,
	logger *zap.Logger,
	logicConfig configs.Logic,
	appArguments utils.Arguments,
) (Judge, error) {
	submissionRetryDelayDuration, err := logicConfig.Judge.GetSubmissionRetryDelayDuration()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get submission retry delay duration")
		return nil, err
	}

	j := &judge{
		problemDataAccessor:             problemDataAccessor,
		submissionDataAccessor:          submissionDataAccessor,
		testCaseDataAccessor:            testCaseDataAccessor,
		problemTestCaseHashDataAccessor: problemTestCaseHashDataAccessor,
		db:                              db,
		apiClient:                       apiClient,
		logger:                          logger,
		logicConfig:                     logicConfig,
		appArguments:                    appArguments,
		workerPool:                      workerpool.New(1),
		languageToCompileLogic:          make(map[string]Compile),
		languageToTestCaseRunLogic:      make(map[string]TestCaseRun),
		submissionRetryDelayDuration:    submissionRetryDelayDuration,
	}

	if appArguments.Distributed && !appArguments.Worker {
		// Skip pulling of Docker image on Distributed host because we are not judging there
		return j, nil
	}

	for _, languageConfig := range logicConfig.Judge.Languages {
		compile, compileErr := NewCompile(dockerClient, logger, languageConfig.Value, languageConfig.Compile)
		if compileErr != nil {
			return nil, compileErr
		}

		j.languageToCompileLogic[languageConfig.Value] = compile

		testCaseRun, testCaseRunErr := NewTestCaseRun(dockerClient, logger, languageConfig.Value, languageConfig.TestCaseRun)
		if testCaseRunErr != nil {
			return nil, testCaseRunErr
		}

		j.languageToTestCaseRunLogic[languageConfig.Value] = testCaseRun
	}

	return j, nil
}

func (j judge) validateProblemHash(ctx context.Context, problem *db.Problem) error {
	logger := utils.LoggerWithContext(ctx, j.logger).With(zap.String("problem_uuid", problem.UUID))

	getProblemResponse, err := j.apiClient.GetProblem(ctx, &rpc.GetProblemRequest{
		UUID: problem.UUID,
	})
	if err != nil {
		return err
	}

	problemTestCaseHash, err := j.problemTestCaseHashDataAccessor.
		GetProblemTestCaseHashOfProblem(ctx, uint64(problem.ID))
	if err != nil {
		return err
	}

	if getProblemResponse.Problem.TestCaseHash != problemTestCaseHash.Hash {
		logger.Error("test case hash of problem is not equal to the latest from host, will not judge")
		return errors.New("test case hash of problem is not equal to the latest from host, will not judge")
	}

	return nil
}

func (j judge) updateSubmissionStatusAndResult(
	ctx context.Context,
	submission *db.Submission,
	status db.SubmissionStatus,
	result db.SubmissionResult,
) error {
	submission.Status = status
	submission.Result = result
	return j.submissionDataAccessor.UpdateSubmission(ctx, submission)
}

func (j judge) judgeDBSubmissionProblemAndTestCase(
	ctx context.Context,
	problem *db.Problem,
	testCase *db.TestCase,
	compileOutput CompileOutput,
	runLogic TestCaseRun,
) (db.SubmissionResult, error) {
	logger := utils.LoggerWithContext(ctx, j.logger).
		With(zap.Uint("problem_id", problem.ID)).
		With(zap.Uint("test_case_id", testCase.ID))
	logger.Info("running submission against test case")

	runOutput, err := runLogic.Run(
		ctx,
		compileOutput.ProgramFilePath,
		testCase.Input,
		problem.TimeLimitInMillisecond,
		problem.MemoryLimitInByte,
	)
	if err != nil {
		return 0, err
	}

	if runOutput.TimeLimitExceeded {
		logger.Info("submission exceeded time limit")
		return db.SubmissionResultTimeLimitExceeded, nil
	}

	if runOutput.MemoryLimitExceeded {
		logger.Info("submission exceeded memory limit")
		return db.SubmissionResultMemoryLimitExceeded, nil
	}

	if runOutput.ReturnCode != 0 {
		logger.With(zap.Int64("return_code", runOutput.ReturnCode)).Info("submission has runtime error")
		return db.SubmissionResultRuntimeError, nil
	}

	if runOutput.StdOut != testCase.Output {
		logger.Info("submission gave incorrect output")
		return db.SubmissionResultWrongAnswer, nil
	}

	return db.SubmissionResultOK, nil
}

func (j judge) JudgeSubmission(
	ctx context.Context,
	problem *db.Problem,
	language string,
	content string,
) (db.SubmissionResult, error) {
	logger := utils.LoggerWithContext(ctx, j.logger).With(zap.String("problemUUID", problem.UUID))

	if j.appArguments.Distributed {
		logger.Info("validating problem hash")
		if err := j.validateProblemHash(ctx, problem); err != nil {
			return 0, err
		}
	}

	logger.Info("compiling submission")
	compileLogic, ok := j.languageToCompileLogic[language]
	if !ok {
		logger.With(zap.String("language", language)).Info("submission has unsupported language")
		return db.SubmissionResultUnsupportedLanguage, nil
	}

	compileOutput, err := compileLogic.Compile(ctx, content)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to compile submission")
		return 0, err
	}

	if compileOutput.ProgramFilePath == "" {
		logger.With(zap.Any("compile_output", compileOutput)).Info("submission has compile error")
		return db.SubmissionResultCompileError, nil
	}

	defer func() {
		err = os.Remove(compileOutput.ProgramFilePath)
		if err != nil {
			logger.
				With(zap.String("program_file_path", compileOutput.ProgramFilePath)).
				With(zap.Error(err)).
				Error("failed to remove program file")
		}
	}()

	logger.Info("running submission against test cases")
	testCaseID, err := j.testCaseDataAccessor.GetTestCaseIDListOfProblem(ctx, uint64(problem.ID))
	if err != nil {
		return 0, err
	}

	runLogic := j.languageToTestCaseRunLogic[language]
	for _, testCaseID := range testCaseID {
		testCase, testCaseErr := j.testCaseDataAccessor.GetTestCase(ctx, testCaseID)
		if testCaseErr != nil {
			return 0, testCaseErr
		}

		result, testCaseErr := j.
			judgeDBSubmissionProblemAndTestCase(ctx, problem, testCase, compileOutput, runLogic)
		if testCaseErr != nil {
			return 0, testCaseErr
		}

		if result != db.SubmissionResultOK {
			return result, nil
		}
	}

	logger.Info("submission passed")
	return db.SubmissionResultOK, nil
}

func (j judge) judgeLocalSubmission(ctx context.Context, submissionID uint64) {
	var (
		logger     = utils.LoggerWithContext(ctx, j.logger).With(zap.Uint64("submission_id", submissionID))
		submission *db.Submission
		err        error
	)

	logger.Info("retrieving submission information")
	if txErr := j.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		submission, err = j.submissionDataAccessor.WithDB(tx).GetSubmission(ctx, submissionID)
		if err != nil {
			return err
		}

		if submission.Status != db.SubmissionStatusSubmitted {
			logger.Error("status of submission is not submitted")
			return errors.New("status of submission is not submitted")
		}

		submission.Status = db.SubmissionStatusExecuting
		err = j.submissionDataAccessor.WithDB(tx).UpdateSubmission(ctx, submission)
		if err != nil {
			return err
		}

		return nil
	}); txErr != nil {
		logger.With(zap.Error(txErr)).Error("failed to retrieve local submission information")
		return
	}

	revertFunc := func(revertErr error) {
		logger.
			With(zap.Error(revertErr)).
			Error("encountered error while judging submission, reverting status to submitted")

		revertErr = j.updateSubmissionStatusAndResult(ctx, submission, db.SubmissionStatusSubmitted, 0)
		if revertErr != nil {
			logger.With(zap.Error(revertErr)).Error("failed to revert submission status to submitted")
		}

		time.AfterFunc(j.submissionRetryDelayDuration, func() {
			j.ScheduleJudgeLocalSubmission(submissionID)
		})
	}

	problem, err := j.problemDataAccessor.GetProblem(ctx, submission.OfProblemID)
	if err != nil {
		return
	}

	result, err := j.JudgeSubmission(ctx, problem, submission.Language, submission.Content)
	if err != nil {
		revertFunc(err)
		return
	}

	err = j.updateSubmissionStatusAndResult(ctx, submission, db.SubmissionStatusFinished, result)
	if err != nil {
		revertFunc(err)
	}
}

func (j judge) ScheduleJudgeLocalSubmission(submissionID uint64) {
	j.workerPool.Submit(func() { j.judgeLocalSubmission(context.Background(), submissionID) })
}

func (j judge) judgeLDistributedSubmission(ctx context.Context) {
	logger := utils.LoggerWithContext(ctx, j.logger)

	response, err := j.apiClient.GetAndUpdateFirstSubmittedSubmissionToExecuting(
		ctx, &rpc.GetAndUpdateFirstSubmittedSubmissionToExecutingRequest{})
	if err != nil {
		return
	}

	problem, err := j.problemDataAccessor.GetProblemByUUID(ctx, response.Submission.Problem.UUID)
	if err != nil {
		return
	}

	if problem == nil {
		logger.
			With(zap.Uint64("submission_id", response.Submission.ID)).
			With(zap.String("problem_uuid", response.Submission.Problem.UUID)).
			Error("problem is not yet synced on worker, will not judge")
		return
	}

	submission := response.Submission
	result, err := j.JudgeSubmission(ctx, problem, submission.Language, submission.Content)
	if err != nil {
		return
	}

	_, err = j.apiClient.UpdateSubmission(ctx, &rpc.UpdateSubmissionRequest{
		ID:     submission.ID,
		Status: uint8(rpc.SubmissionStatusFinished),
		Result: uint8(result),
	})
	if err != nil {
		return
	}
}

func (j judge) ScheduleJudgeDistributedFirstSubmittedSubmission() {
	j.workerPool.Submit(func() { j.judgeLDistributedSubmission(context.Background()) })
}
