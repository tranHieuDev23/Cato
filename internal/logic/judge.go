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
	"github.com/tranHieuDev23/cato/internal/utils"
)

type Judge interface {
	JudgeSubmission(ctx context.Context, id uint64) error
	ScheduleSubmissionToJudge(id uint64)
}

type judge struct {
	problemDataAccessor          db.ProblemDataAccessor
	submissionDataAccessor       db.SubmissionDataAccessor
	testCaseDataAccessor         db.TestCaseDataAccessor
	db                           *gorm.DB
	logger                       *zap.Logger
	logicConfig                  configs.Logic
	isLocal                      bool
	workerPool                   *workerpool.WorkerPool
	languageToCompileLogic       map[string]Compile
	languageToTestCaseRunLogic   map[string]TestCaseRun
	submissionRetryDelayDuration time.Duration
}

func NewJudge(
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	testCaseDataAccessor db.TestCaseDataAccessor,
	dockerClient *client.Client,
	db *gorm.DB,
	logger *zap.Logger,
	logicConfig configs.Logic,
	isLocal bool,
) (Judge, error) {
	submissionRetryDelayDuration, err := logicConfig.Judge.GetSubmissionRetryDelayDuration()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get submission retry delay duration")
		return nil, err
	}

	j := &judge{
		problemDataAccessor:          problemDataAccessor,
		submissionDataAccessor:       submissionDataAccessor,
		testCaseDataAccessor:         testCaseDataAccessor,
		db:                           db,
		logger:                       logger,
		logicConfig:                  logicConfig,
		isLocal:                      isLocal,
		workerPool:                   workerpool.New(1),
		languageToCompileLogic:       make(map[string]Compile),
		languageToTestCaseRunLogic:   make(map[string]TestCaseRun),
		submissionRetryDelayDuration: submissionRetryDelayDuration,
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

func (j judge) validateProblemHash(_ context.Context, _ *db.Problem) error {
	// TODO: Fill this in when implement worker APIs
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
	submission *db.Submission,
	problem *db.Problem,
	testCase *db.TestCase,
	compileOutput CompileOutput,
	runLogic TestCaseRun,
) (bool, error) {
	logger := utils.LoggerWithContext(ctx, j.logger).
		With(zap.Uint("submission_id", submission.ID)).
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
		return false, err
	}

	if runOutput.TimeLimitExceeded {
		logger.Info("submission exceeded time limit")
		err = j.updateSubmissionStatusAndResult(
			ctx, submission, db.SubmissionStatusFinished, db.SubmissionResultTimeLimitExceeded)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	if runOutput.MemoryLimitExceeded {
		logger.Info("submission exceeded memory limit")
		err = j.updateSubmissionStatusAndResult(
			ctx, submission, db.SubmissionStatusFinished, db.SubmissionResultMemoryLimitExceeded)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	if runOutput.ReturnCode != 0 {
		logger.With(zap.Int64("return_code", runOutput.ReturnCode)).Info("submission has runtime error")
		err = j.updateSubmissionStatusAndResult(
			ctx, submission, db.SubmissionStatusFinished, db.SubmissionResultRuntimeError)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	if runOutput.StdOut != testCase.Output {
		logger.Info("submission gave incorrect output")
		err = j.updateSubmissionStatusAndResult(
			ctx, submission, db.SubmissionStatusFinished, db.SubmissionResultWrongAnswer)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func (j judge) judgeDBSubmission(ctx context.Context, submission *db.Submission) error {
	logger := utils.LoggerWithContext(ctx, j.logger).With(zap.Uint("submission_id", submission.ID))

	problem, err := j.problemDataAccessor.GetProblem(ctx, submission.OfProblemID)
	if err != nil {
		return err
	}

	if problem == nil {
		logger.With(zap.Uint64("problem_id", submission.OfProblemID)).Error("cannot find problem")
		return errors.New("cannot find problem")
	}

	if !j.isLocal {
		logger.Info("validating problem hash")
		err = j.validateProblemHash(ctx, problem)
		if err != nil {
			return err
		}
	}

	logger.Info("retrieving test case information")
	testCaseID, err := j.testCaseDataAccessor.GetTestCaseIDListOfProblem(ctx, uint64(problem.ID))
	if err != nil {
		return err
	}

	logger.Info("compiling submission")
	compileLogic, ok := j.languageToCompileLogic[submission.Language]
	if !ok {
		logger.With(zap.String("language", submission.Language)).Info("submission has unsupported language")
		return j.updateSubmissionStatusAndResult(
			ctx, submission, db.SubmissionStatusFinished, db.SubmissionResultUnsupportedLanguage)
	}

	compileOutput, err := compileLogic.Compile(ctx, submission.Content)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to compile submission")
		return err
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

	if compileOutput.ProgramFilePath == "" {
		logger.With(zap.Any("compile_output", compileOutput)).Info("submission has compile error")
		return j.updateSubmissionStatusAndResult(
			ctx, submission, db.SubmissionStatusFinished, db.SubmissionResultCompileError)
	}

	logger.Info("running submission against test cases")
	runLogic := j.languageToTestCaseRunLogic[submission.Language]
	for _, testCaseID := range testCaseID {
		testCase, testCaseErr := j.testCaseDataAccessor.GetTestCase(ctx, testCaseID)
		if testCaseErr != nil {
			return testCaseErr
		}

		passed, testCaseErr := j.
			judgeDBSubmissionProblemAndTestCase(ctx, submission, problem, testCase, compileOutput, runLogic)
		if testCaseErr != nil {
			return testCaseErr
		}

		if !passed {
			return nil
		}
	}

	logger.Info("submission passed")
	return j.updateSubmissionStatusAndResult(ctx, submission, db.SubmissionStatusFinished, db.SubmissionResultOK)
}

func (j judge) JudgeSubmission(ctx context.Context, id uint64) error {
	var (
		logger     = utils.LoggerWithContext(ctx, j.logger).With(zap.Uint64("id", id))
		submission *db.Submission
		err        error
	)

	logger.Info("retrieving submission information")
	if txErr := j.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		submission, err = j.submissionDataAccessor.WithDB(tx).GetSubmission(ctx, id)
		if err != nil {
			return err
		}

		if submission == nil {
			logger.Error("cannot find submission")
			return errors.New("cannot find submission")
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
		return txErr
	}

	err = j.judgeDBSubmission(ctx, submission)
	if err != nil {
		logger.With(zap.Error(err)).Error("encountered error while judging submission, reverting status to submitted")

		if revertErr := j.updateSubmissionStatusAndResult(
			ctx, submission, db.SubmissionStatusSubmitted, 0,
		); revertErr != nil {
			logger.With(zap.Error(revertErr)).Error("failed to revert submission status to submitted")
		}

		time.AfterFunc(j.submissionRetryDelayDuration, func() {
			j.ScheduleSubmissionToJudge(id)
		})

		return err
	}

	return nil
}

func (j judge) ScheduleSubmissionToJudge(id uint64) {
	j.workerPool.Submit(func() {
		if err := j.JudgeSubmission(context.Background(), id); err != nil {
			j.logger.With(zap.Error(err)).Error("error occurred when judging submission")
		}
	})
}

type LocalJudge Judge

func NewLocalJudge(
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	testCaseDataAccessor db.TestCaseDataAccessor,
	dockerClient *client.Client,
	db *gorm.DB,
	logger *zap.Logger,
	logicConfig configs.Logic,
) (LocalJudge, error) {
	return NewJudge(
		problemDataAccessor,
		submissionDataAccessor,
		testCaseDataAccessor,
		dockerClient,
		db,
		logger,
		logicConfig,
		true,
	)
}

type DistributedJudge Judge

func NewDistributedJudge(
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	testCaseDataAccessor db.TestCaseDataAccessor,
	dockerClient *client.Client,
	db *gorm.DB,
	logger *zap.Logger,
	logicConfig configs.Logic,
) (DistributedJudge, error) {
	return NewJudge(
		problemDataAccessor,
		submissionDataAccessor,
		testCaseDataAccessor,
		dockerClient,
		db,
		logger,
		logicConfig,
		false,
	)
}
