package logic

import (
	"context"
	"errors"

	"github.com/docker/docker/client"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type Judge interface {
	JudgeSubmission(ctx context.Context, id uint64) error
}

type judge struct {
	testCaseLogic              TestCase
	problemDataAccessor        db.ProblemDataAccessor
	submissionDataAccessor     db.SubmissionDataAccessor
	testCaseDataAccessor       db.TestCaseDataAccessor
	dockerClient               *client.Client
	db                         *gorm.DB
	logger                     *zap.Logger
	logicConfig                configs.Logic
	shouldValidateProblemHash  bool
	languageToCompileLogic     map[string]Compile
	languageToTestCaseRunLogic map[string]TestCaseRun
}

func NewJudge(
	testCaseLogic TestCase,
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	testCaseDataAccessor db.TestCaseDataAccessor,
	problemTestCaseHashDataAccessor db.ProblemTestCaseHashDataAccessor,
	dockerClient *client.Client,
	db *gorm.DB,
	logger *zap.Logger,
	logicConfig configs.Logic,
	shouldValidateProblemHash bool,
) (Judge, error) {
	e := &judge{
		problemDataAccessor:        problemDataAccessor,
		submissionDataAccessor:     submissionDataAccessor,
		testCaseDataAccessor:       testCaseDataAccessor,
		db:                         db,
		logger:                     logger,
		logicConfig:                logicConfig,
		shouldValidateProblemHash:  shouldValidateProblemHash,
		languageToCompileLogic:     make(map[string]Compile),
		languageToTestCaseRunLogic: make(map[string]TestCaseRun),
	}

	for language, config := range logicConfig.Judge.Languages {
		compile, err := NewCompile(dockerClient, logger, language, config.Compile)
		if err != nil {
			return nil, err
		}

		e.languageToCompileLogic[language] = compile

		testCaseRun, err := NewTestCaseRun(dockerClient, logger, language, config.TestCaseRun)
		if err != nil {
			return nil, err
		}

		e.languageToTestCaseRunLogic[language] = testCaseRun
	}

	return e, nil
}

func (e judge) validateProblemHash(ctx context.Context, problem *db.Problem) error {
	// TODO: Fill this in when implement worker APIs
	return nil
}

func (e judge) updateSubmissionStatusAndResult(
	ctx context.Context,
	submission *db.Submission,
	status db.SubmissionStatus,
	result db.SubmissionResult,
) error {
	submission.Status = status
	submission.Result = result
	return e.submissionDataAccessor.UpdateSubmission(ctx, submission)
}

func (e judge) judgeDBSubmission(ctx context.Context, submission *db.Submission) error {
	logger := utils.LoggerWithContext(ctx, e.logger).With(zap.Uint("id", submission.ID))

	problem, err := e.problemDataAccessor.GetProblem(ctx, submission.OfProblemID)
	if err != nil {
		return err
	}

	if problem == nil {
		logger.With(zap.Uint64("problem_id", submission.OfProblemID)).Error("cannot find problem")
		return errors.New("cannot find problem")
	}

	if e.shouldValidateProblemHash {
		logger.Info("validating problem hash")
		if err := e.validateProblemHash(ctx, problem); err != nil {
			return err
		}
	}

	logger.Info("retrieving test case information")
	testCaseID, err := e.testCaseDataAccessor.GetTestCaseIDListOfProblem(ctx, uint64(problem.ID))
	if err != nil {
		return err
	}

	logger.Info("compiling submission")
	compileLogic, ok := e.languageToCompileLogic[submission.Language]
	if !ok {
		logger.With(zap.String("language", submission.Language)).Info("submission has supported language")
		return e.updateSubmissionStatusAndResult(
			ctx, submission, db.SubmissionStatusFinished, db.SubmissionResultUnsupportedLanguage)
	}

	compileOutput, err := compileLogic.Compile(ctx, submission.Content)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to compile submission")
		return err
	}

	if compileOutput.ProgramFilePath == "" {
		logger.With(zap.Any("compile_output", compileOutput)).Info("submission has compile error")
		return e.updateSubmissionStatusAndResult(
			ctx, submission, db.SubmissionStatusFinished, db.SubmissionResultCompileError)
	}

	logger.Info("running submission against test cases")
	runLogic := e.languageToTestCaseRunLogic[submission.Language]
	for _, testCaseID := range testCaseID {
		testCase, err := e.testCaseDataAccessor.GetTestCase(ctx, testCaseID)
		if err != nil {
			return err
		}

		logger.With(zap.Uint64("test_case_id", testCaseID)).Info("running submission against test case")
		runOutput, err := runLogic.Run(
			ctx,
			compileOutput.ProgramFilePath,
			testCase.Input,
			problem.TimeLimitInMillisecond,
			problem.MemoryLimitInByte,
		)
		if err != nil {
			return err
		}

		if runOutput.ReturnCode != 0 {
			logger.
				With(zap.Uint64("test_case_id", testCaseID)).
				With(zap.Int64("return_code", runOutput.ReturnCode)).
				Info("submission has runtime error")
			return e.updateSubmissionStatusAndResult(
				ctx, submission, db.SubmissionStatusFinished, db.SubmissionResultRuntimeError)
		}

		if runOutput.StdOut != testCase.Output {
			logger.With(zap.Uint64("test_case_id", testCaseID)).Info("submission gave incorrect output")
			return e.updateSubmissionStatusAndResult(
				ctx, submission, db.SubmissionStatusFinished, db.SubmissionResultWrongAnswer)
		}
	}

	logger.Info("submission passed")
	return e.updateSubmissionStatusAndResult(ctx, submission, db.SubmissionStatusFinished, db.SubmissionResultOK)
}

func (e judge) JudgeSubmission(ctx context.Context, id uint64) error {
	var (
		logger     = utils.LoggerWithContext(ctx, e.logger).With(zap.Uint64("id", id))
		submission *db.Submission
		err        error
	)

	logger.Info("retrieving submission information")
	if txErr := e.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		submission, err = e.submissionDataAccessor.WithDB(tx).GetSubmission(ctx, id)
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
		if err := e.submissionDataAccessor.WithDB(tx).UpdateSubmission(ctx, submission); err != nil {
			return err
		}

		return nil
	}); txErr != nil {
		return txErr
	}

	if err := e.judgeDBSubmission(ctx, submission); err != nil {
		logger.With(zap.Error(err)).Error("encountered error while judging submission, reverting status to submitted")

		if revertErr := e.updateSubmissionStatusAndResult(
			ctx, submission, db.SubmissionStatusSubmitted, 0,
		); revertErr != nil {
			logger.With(zap.Error(revertErr)).Error("failed to revert submission status to submitted")
		}

		return err
	}

	return nil
}

type LocalJudge Judge

func NewLocalJudge(
	testCaseLogic TestCase,
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	testCaseDataAccessor db.TestCaseDataAccessor,
	problemTestCaseHashDataAccessor db.ProblemTestCaseHashDataAccessor,
	dockerClient *client.Client,
	db *gorm.DB,
	logger *zap.Logger,
	logicConfig configs.Logic,
) (LocalJudge, error) {
	return NewJudge(
		testCaseLogic,
		problemDataAccessor,
		submissionDataAccessor,
		testCaseDataAccessor,
		problemTestCaseHashDataAccessor,
		dockerClient,
		db,
		logger,
		logicConfig,
		false,
	)
}

type RemoteJudge Judge

func NewRemoteJudge(
	testCaseLogic TestCase,
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	testCaseDataAccessor db.TestCaseDataAccessor,
	problemTestCaseHashDataAccessor db.ProblemTestCaseHashDataAccessor,
	dockerClient *client.Client,
	db *gorm.DB,
	logger *zap.Logger,
	logicConfig configs.Logic,
) (LocalJudge, error) {
	return NewJudge(
		testCaseLogic,
		problemDataAccessor,
		submissionDataAccessor,
		testCaseDataAccessor,
		problemTestCaseHashDataAccessor,
		dockerClient,
		db,
		logger,
		logicConfig,
		false,
	)
}
