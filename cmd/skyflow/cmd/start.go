package cmd

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/goodaye/wire"
	"github.com/skyflow-workflow/skyflow_backbend/server/apiserver"
	"github.com/skyflow-workflow/skyflow_backbend/workflow"
	"github.com/spf13/cobra"
)

var (
	trpc_conf    string = "./trpc_go.yaml" // Configuration file for the TRPC framework
	skyflow_conf string = "./skyflow.yaml" // Configuration file for the skyflow server
)

func init() {
	startCmd.Flags().StringVarP(&trpc_conf, "trpc_config", "c", "./trpc_go.yaml", " trpc configuration file, default is ./trpc_go.yaml")
	startCmd.Flags().StringVarP(&skyflow_conf, "skyflow_config", "", "./skyflow.yaml", " skyflow configuration file, default is ./skyflow.yaml")

}

var startCmd = &cobra.Command{
	Use:   "start [flags] [args]",
	Short: "Start the skyflow workflow engine",
	Long: `Start the skyflow workflow engine, which allows you to run workflows defined in JSON format.
args:
	api - start the API server
	dispatcher - start the dispatcher
`,
	Run: StartCommand,
}

func StartCommand(cmd *cobra.Command, args []string) {
	var err error
	// initialize trpc server
	err = InitializeTrpc(trpc_conf)
	if err != nil {
		slog.Error("Error initializing trpc server", "error", err)
		return
	}

	sfConfig, err := LoadConfig(skyflow_conf)
	if err != nil {
		slog.Error("Error loading skyflow configuration", "error", err)
		return
	}
	wfSvc, err := LoadServices(sfConfig)
	if err != nil {
		slog.Error("Error loading services", "error", err)
		return
	}

	wg := sync.WaitGroup{}

	for _, arg := range args {
		switch arg {
		case "api":
			wg.Add(1)
			go func() {
				defer wg.Done()
				slog.Info("Starting API server...")
				StartAPIServer(wfSvc)
			}()
		case "dispatcher":
			wg.Add(1)
			go func() {
				defer wg.Done()
				slog.Info("Starting Dispatcher...")
				StartDispatcher(wfSvc)
			}()
		default:
			err = fmt.Errorf("unknown command:  %s", arg)
		}
		if err != nil {
			slog.Error("Error starting skyflow", "error", err)
			return
		}
	}

	// 自动会调用 wire.Stop
	err = wire.Run()
	if err != nil {
		slog.Error("Error running wire", "error", err)
	}
	// receive  stop signal
	slog.Info("stop signal received")
	wg.Wait()

	slog.Info("stop skyflow  finished")

}

func StartAPIServer(wfSvc workflow.WorkflowService) {
	// Initialize the API server
	apiServer := apiserver.NewAPIServer(wfSvc)
	// Start the API server
	apiServer.Start()
}

func StartDispatcher(wfSvc workflow.WorkflowService) {
	// // Initialize the dispatcher
	// dispatcher := dispatcher.NewDispatcher(trpc_conf)

	// // Start the dispatcher
	// if err := dispatcher.Start(); err != nil {
	// 	fmt.Println("Failed to start dispatcher:", err)
	// }
}
