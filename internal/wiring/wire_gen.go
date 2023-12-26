// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wiring

import (
	"github.com/google/wire"

	"github.com/tranHieuDev23/cato/internal/app"
	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess"
	"github.com/tranHieuDev23/cato/internal/dataaccess/cato"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
	"github.com/tranHieuDev23/cato/internal/handlers"
	"github.com/tranHieuDev23/cato/internal/handlers/http"
	"github.com/tranHieuDev23/cato/internal/handlers/http/middlewares"
	"github.com/tranHieuDev23/cato/internal/handlers/jobs"
	"github.com/tranHieuDev23/cato/internal/logic"
	"github.com/tranHieuDev23/cato/internal/utils"
)

// Injectors from wire.go:

func InitializeLocalCato(filePath configs.ConfigFilePath, args utils.Arguments) (*app.LocalCato, func(), error) {
	logger, cleanup, err := utils.InitializeLogger()
	if err != nil {
		return nil, nil, err
	}
	config, err := configs.NewConfig(filePath)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	database := config.Database
	gormDB, err := db.InitializeDB(logger, database)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	migrator := db.NewMigrator(gormDB, logger)
	auth := config.Auth
	hash := auth.Hash
	logicHash := logic.NewHash(hash, logger)
	accountDataAccessor := db.NewAccountDataAccessor(gormDB, logger)
	token := auth.Token
	logicToken, err := logic.NewToken(accountDataAccessor, token, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	role := logic.NewRole(logger)
	accountPasswordDataAccessor := db.NewAccountPasswordDataAccessor(gormDB, logger)
	configsLogic := config.Logic
	account := logic.NewAccount(logicHash, logicToken, role, accountDataAccessor, accountPasswordDataAccessor, gormDB, logger, configsLogic, args)
	createFirstAccounts := jobs.NewCreateFirstAccounts(account)
	problemDataAccessor := db.NewProblemDataAccessor(gormDB, logger)
	submissionDataAccessor := db.NewSubmissionDataAccessor(gormDB, logger)
	testCaseDataAccessor := db.NewTestCaseDataAccessor(gormDB, logger)
	client, err := utils.InitializeDockerClient()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	httpClientWithAuthToken, err := cato.NewHTTPClientWithAuthToken(args, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	apiClient := cato.InitializeAuthenticatedClient(args, httpClientWithAuthToken)
	judge, err := logic.NewJudge(problemDataAccessor, submissionDataAccessor, testCaseDataAccessor, client, gormDB, apiClient, logger, configsLogic, args)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	submission := logic.NewSubmission(logicToken, role, judge, accountDataAccessor, problemDataAccessor, submissionDataAccessor, gormDB, logger, args)
	scheduleSubmittedExecutingSubmissionToJudge := jobs.NewScheduleSubmittedExecutingSubmissionToJudge(submission)
	problemTestCaseHashDataAccessor := db.NewProblemTestCaseHashDataAccessor(gormDB, logger)
	testCase := logic.NewTestCase(logicToken, role, problemDataAccessor, testCaseDataAccessor, problemTestCaseHashDataAccessor, gormDB, logger, configsLogic)
	problemExampleDataAccessor := db.NewProblemExampleDataAccessor(gormDB, logger)
	problem := logic.NewProblem(logicToken, role, testCase, accountDataAccessor, problemDataAccessor, problemExampleDataAccessor, problemTestCaseHashDataAccessor, testCaseDataAccessor, submissionDataAccessor, logger, gormDB)
	apiServer := http.NewAPIServerHandler(account, problem, testCase, submission, configsLogic, logger, args)
	v := middlewares.InitializePJRPCMiddlewareList()
	httpAuth, err := middlewares.NewHTTPAuth(logicToken, token, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	v2 := middlewares.InitalizeHTTPMiddlewareList(httpAuth)
	spaHandler := http.NewSPAHandler()
	configsHTTP := config.HTTP
	server := http.NewServer(apiServer, v, v2, spaHandler, logger, configsHTTP, args)
	localCato := app.NewLocalCato(migrator, createFirstAccounts, scheduleSubmittedExecutingSubmissionToJudge, server, logger)
	return localCato, func() {
		cleanup()
	}, nil
}

func InitializeDistributedHostCato(filePath configs.ConfigFilePath, args utils.Arguments) (*app.DistributedHostCato, func(), error) {
	logger, cleanup, err := utils.InitializeLogger()
	if err != nil {
		return nil, nil, err
	}
	config, err := configs.NewConfig(filePath)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	database := config.Database
	gormDB, err := db.InitializeDB(logger, database)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	migrator := db.NewMigrator(gormDB, logger)
	auth := config.Auth
	hash := auth.Hash
	logicHash := logic.NewHash(hash, logger)
	accountDataAccessor := db.NewAccountDataAccessor(gormDB, logger)
	token := auth.Token
	logicToken, err := logic.NewToken(accountDataAccessor, token, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	role := logic.NewRole(logger)
	accountPasswordDataAccessor := db.NewAccountPasswordDataAccessor(gormDB, logger)
	configsLogic := config.Logic
	account := logic.NewAccount(logicHash, logicToken, role, accountDataAccessor, accountPasswordDataAccessor, gormDB, logger, configsLogic, args)
	createFirstAccounts := jobs.NewCreateFirstAccounts(account)
	problemDataAccessor := db.NewProblemDataAccessor(gormDB, logger)
	testCaseDataAccessor := db.NewTestCaseDataAccessor(gormDB, logger)
	problemTestCaseHashDataAccessor := db.NewProblemTestCaseHashDataAccessor(gormDB, logger)
	testCase := logic.NewTestCase(logicToken, role, problemDataAccessor, testCaseDataAccessor, problemTestCaseHashDataAccessor, gormDB, logger, configsLogic)
	problemExampleDataAccessor := db.NewProblemExampleDataAccessor(gormDB, logger)
	submissionDataAccessor := db.NewSubmissionDataAccessor(gormDB, logger)
	problem := logic.NewProblem(logicToken, role, testCase, accountDataAccessor, problemDataAccessor, problemExampleDataAccessor, problemTestCaseHashDataAccessor, testCaseDataAccessor, submissionDataAccessor, logger, gormDB)
	client, err := utils.InitializeDockerClient()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	httpClientWithAuthToken, err := cato.NewHTTPClientWithAuthToken(args, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	apiClient := cato.InitializeAuthenticatedClient(args, httpClientWithAuthToken)
	judge, err := logic.NewJudge(problemDataAccessor, submissionDataAccessor, testCaseDataAccessor, client, gormDB, apiClient, logger, configsLogic, args)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	submission := logic.NewSubmission(logicToken, role, judge, accountDataAccessor, problemDataAccessor, submissionDataAccessor, gormDB, logger, args)
	apiServer := http.NewAPIServerHandler(account, problem, testCase, submission, configsLogic, logger, args)
	v := middlewares.InitializePJRPCMiddlewareList()
	httpAuth, err := middlewares.NewHTTPAuth(logicToken, token, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	v2 := middlewares.InitalizeHTTPMiddlewareList(httpAuth)
	spaHandler := http.NewSPAHandler()
	configsHTTP := config.HTTP
	server := http.NewServer(apiServer, v, v2, spaHandler, logger, configsHTTP, args)
	distributedHostCato := app.NewDistributedHostCato(migrator, createFirstAccounts, server, logger)
	return distributedHostCato, func() {
		cleanup()
	}, nil
}

func InitializeDistributedWorkerCato(filePath configs.ConfigFilePath, args utils.Arguments) (*app.DistributedWorkerCato, func(), error) {
	logger, cleanup, err := utils.InitializeLogger()
	if err != nil {
		return nil, nil, err
	}
	config, err := configs.NewConfig(filePath)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	database := config.Database
	gormDB, err := db.InitializeDB(logger, database)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	migrator := db.NewMigrator(gormDB, logger)
	accountDataAccessor := db.NewAccountDataAccessor(gormDB, logger)
	auth := config.Auth
	token := auth.Token
	logicToken, err := logic.NewToken(accountDataAccessor, token, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	role := logic.NewRole(logger)
	problemDataAccessor := db.NewProblemDataAccessor(gormDB, logger)
	submissionDataAccessor := db.NewSubmissionDataAccessor(gormDB, logger)
	testCaseDataAccessor := db.NewTestCaseDataAccessor(gormDB, logger)
	client, err := utils.InitializeDockerClient()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	httpClientWithAuthToken, err := cato.NewHTTPClientWithAuthToken(args, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	apiClient := cato.InitializeAuthenticatedClient(args, httpClientWithAuthToken)
	configsLogic := config.Logic
	judge, err := logic.NewJudge(problemDataAccessor, submissionDataAccessor, testCaseDataAccessor, client, gormDB, apiClient, logger, configsLogic, args)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	submission := logic.NewSubmission(logicToken, role, judge, accountDataAccessor, problemDataAccessor, submissionDataAccessor, gormDB, logger, args)
	scheduleSubmittedExecutingSubmissionToJudge := jobs.NewScheduleSubmittedExecutingSubmissionToJudge(submission)
	distributedWorkerCato := app.NewDistributedWorkerCato(migrator, scheduleSubmittedExecutingSubmissionToJudge, logger)
	return distributedWorkerCato, func() {
		cleanup()
	}, nil
}

// wire.go:

var WireSet = wire.NewSet(utils.WireSet, configs.WireSet, dataaccess.WireSet, logic.WireSet, handlers.WireSet, app.WireSet)
