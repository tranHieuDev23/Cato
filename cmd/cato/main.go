package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/tranHieuDev23/cato/internal/app"
	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/wiring"
)

const (
	flagConfigFilePath = "config-file-path"
	flagDistributed    = "distributed"
	flagWorker         = "worker"
)

func main() {
	rootCommand := &cobra.Command{
		Use: "Start the program",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFilePath, err := cmd.Flags().GetString(flagConfigFilePath)
			if err != nil {
				return err
			}

			distributed, err := cmd.Flags().GetBool(flagDistributed)
			if err != nil {
				return err
			}

			worker, err := cmd.Flags().GetBool(flagWorker)
			if err != nil {
				return err
			}

			var (
				app     app.App
				cleanup func()
			)

			if distributed {
				if worker {
					app, cleanup, err = wiring.InitializeDistributedWorkerCato(configs.ConfigFilePath(configFilePath))
				} else {
					app, cleanup, err = wiring.InitializeDistributedHostCato(configs.ConfigFilePath(configFilePath))
				}
			} else {
				app, cleanup, err = wiring.InitializeLocalCato(configs.ConfigFilePath(configFilePath))
			}

			if err != nil {
				return err
			}

			defer cleanup()

			return app.Start()
		},
	}

	rootCommand.Flags().Bool(
		flagDistributed, false, "If provided, will start the problem in distributed mode.")
	rootCommand.Flags().Bool(
		flagWorker, false, "If provided and --distributed is set, will start the problem in distributed mode as a worker process.")
	rootCommand.Flags().String(
		flagConfigFilePath, "", "If provided, will use the provided config file.")

	if err := rootCommand.Execute(); err != nil {
		log.Panic(err)
	}
}
