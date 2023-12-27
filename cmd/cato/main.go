package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/tranHieuDev23/cato/internal/app"
	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/utils"
	"github.com/tranHieuDev23/cato/internal/wiring"
)

var (
	version    string
	commitHash string
)

const (
	flagDistributed           = "distributed"
	flagWorker                = "worker"
	flagNoBrowser             = "no-browser"
	flagHostAddress           = "host-address"
	flagWorkerAccountName     = "worker-account-name"
	flagWorkerAccountPassword = "worker-account-password"
	flagConfigFilePath        = "config-file-path"
	flagPullImageAtStartUp    = "pull-image-at-startup"
)

func getArguments(cmd *cobra.Command) (utils.Arguments, error) {
	distributed, err := cmd.Flags().GetBool(flagDistributed)
	if err != nil {
		return utils.Arguments{}, err
	}

	worker, err := cmd.Flags().GetBool(flagWorker)
	if err != nil {
		return utils.Arguments{}, err
	}

	noBrowser, err := cmd.Flags().GetBool(flagNoBrowser)
	if err != nil {
		return utils.Arguments{}, err
	}

	hostAddress, err := cmd.Flags().GetString(flagHostAddress)
	if err != nil {
		return utils.Arguments{}, err
	}

	workerAccountName, err := cmd.Flags().GetString(flagWorkerAccountName)
	if err != nil {
		return utils.Arguments{}, err
	}

	workerAccountPassword, err := cmd.Flags().GetString(flagWorkerAccountPassword)
	if err != nil {
		return utils.Arguments{}, err
	}

	pullImageAtStartUp, err := cmd.Flags().GetBool(flagPullImageAtStartUp)
	if err != nil {
		return utils.Arguments{}, err
	}

	return utils.Arguments{
		Distributed:           distributed,
		Worker:                worker,
		NoBrowser:             noBrowser,
		HostAddress:           hostAddress,
		WorkerAccountName:     workerAccountName,
		WorkerAccountPassword: workerAccountPassword,
		PullImageAtStartUp:    pullImageAtStartUp,
	}, nil
}

func main() {
	rootCommand := &cobra.Command{
		Version: fmt.Sprintf("%s-%s", version, commitHash),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				app     app.App
				cleanup func()
			)

			configFilePath, err := cmd.Flags().GetString(flagConfigFilePath)
			if err != nil {
				return err
			}

			arguments, err := getArguments(cmd)
			if err != nil {
				return err
			}

			if arguments.Distributed && arguments.Worker {
				app, cleanup, err = wiring.InitializeWorker(configs.ConfigFilePath(configFilePath), arguments)
			} else {
				app, cleanup, err = wiring.InitializeHost(configs.ConfigFilePath(configFilePath), arguments)
			}

			if err != nil {
				return err
			}

			defer cleanup()

			return app.Start()
		},
	}

	rootCommand.Flags().Bool(
		flagDistributed,
		false,
		"If provided, will start the problem in distributed mode.",
	)
	rootCommand.Flags().Bool(
		flagWorker,
		false,
		"If provided and --distributed is set, will start the program in distributed mode as a worker process.",
	)
	rootCommand.Flags().Bool(
		flagNoBrowser,
		false,
		"If provided, will not open a browser windows when the server starts.",
	)
	rootCommand.Flags().String(
		flagHostAddress,
		"http://127.0.0.1:8080",
		"The address of the host server when running in worker mode.",
	)
	rootCommand.Flags().String(
		flagWorkerAccountName,
		"worker",
		"The worker account name when running in worker mode.",
	)
	rootCommand.Flags().String(
		flagWorkerAccountPassword,
		"changeme",
		"The worker account name when running in worker mode.",
	)
	rootCommand.Flags().String(
		flagConfigFilePath,
		"",
		"If provided, will use the provided config file.",
	)
	rootCommand.Flags().Bool(
		flagPullImageAtStartUp,
		true,
		"Whether to pull Docker images necessary for compiling and executing test case at startup. "+
			"If set to true and Docker fails to pull any of the provided image, the program will exit with non-zero "+
			"error code. ",
	)

	if err := rootCommand.Execute(); err != nil {
		log.Panic(err)
	}
}
