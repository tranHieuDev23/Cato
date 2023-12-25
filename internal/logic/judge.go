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
	problemDataAccessor    db.ProblemDataAccessor
	submissionDataAccessor db.SubmissionDataAccessor
	testCaseDataAccessor   db.TestCaseDataAccessor
	dockerClient           *client.Client
	db                     *gorm.DB
	logger                 *zap.Logger
	logicConfig            configs.Logic

	languageToCompileLogic     map[string]Compile
	languageToTestCaseRunLogic map[string]TestCaseRun
}

func NewJudge(
	problemDataAccessor db.ProblemDataAccessor,
	submissionDataAccessor db.SubmissionDataAccessor,
	testCaseDataAccessor db.TestCaseDataAccessor,
	dockerClient *client.Client,
	db *gorm.DB,
	logger *zap.Logger,
	logicConfig configs.Logic,
) (Judge, error) {
	e := &judge{
		problemDataAccessor:        problemDataAccessor,
		submissionDataAccessor:     submissionDataAccessor,
		testCaseDataAccessor:       testCaseDataAccessor,
		db:                         db,
		logger:                     logger,
		logicConfig:                logicConfig,
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

func (e judge) updateFinishedSubmissionResult(ctx context.Context, submission *db.Submission, result db.SubmissionResult) error {
	submission.Status = db.SubmissionStatusFinished
	submission.Result = result
	return e.submissionDataAccessor.UpdateSubmission(ctx, submission)
}

func (e judge) JudgeSubmission(ctx context.Context, id uint64) error {
	logger := utils.LoggerWithContext(ctx, e.logger).With(zap.Uint64("id", id))

	logger.Info("retrieving submission information")
	submission, err := e.submissionDataAccessor.GetSubmission(ctx, id)
	if err != nil {
		return err
	}

	if submission == nil {
		logger.Error("cannot find submission")
		return errors.New("cannot find submission")
	}

	problem, err := e.problemDataAccessor.GetProblem(ctx, submission.OfProblemID)
	if err != nil {
		return err
	}

	if problem == nil {
		logger.With(zap.Uint64("problem_id", submission.OfProblemID)).Error("cannot find problem")
		return errors.New("cannot find problem")
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
		return e.updateFinishedSubmissionResult(ctx, submission, db.SubmissionResultUnsupportedLanguage)
	}

	compileOutput, err := compileLogic.Compile(ctx, submission.Content)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to compile submission")
		return err
	}

	if compileOutput.ProgramFilePath == "" {
		logger.With(zap.Any("compile_output", compileOutput)).Info("submission has compile error")
		return e.updateFinishedSubmissionResult(ctx, submission, db.SubmissionResultCompileError)
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
			return e.updateFinishedSubmissionResult(ctx, submission, db.SubmissionResultRuntimeError)
		}

		if runOutput.StdOut != testCase.Output {
			logger.With(zap.Uint64("test_case_id", testCaseID)).Info("submission gave incorrect output")
			return e.updateFinishedSubmissionResult(ctx, submission, db.SubmissionResultWrongAnswer)
		}
	}

	logger.Info("submission passed")
	return e.updateFinishedSubmissionResult(ctx, submission, db.SubmissionResultOK)
}
