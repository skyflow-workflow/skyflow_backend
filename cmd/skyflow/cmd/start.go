package cmd

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/goodaye/wire"
	"github.com/skyflow-workflow/skyflow_backbend/server/apiserver"
	"github.com/skyflow-workflow/skyflow_backbend/workflow"
	"github.com/spf13/cobra"
	"trpc.group/trpc-go/trpc-go/server"
)

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
	// initialize trpc server,
	// this will load the configuration file and create a new server instance
	// many plugins depends on trpc server configuration file
	trpcServer := InitializeTrpcSever(trpc_conf)
	sfConfig, err := LoadConfig(skyflow_conf)
	if err != nil {
		slog.Error("Error loading skyflow configuration", "error", err)
		return
	}
	wfSvc, err := LoadService(sfConfig)
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
				StartAPIServer(trpcServer, wfSvc)
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

func StartAPIServer(server *server.Server, wfSvc workflow.WorkflowService) {
	// Initialize the API server
	apiServer := apiserver.NewAPIServer(server, wfSvc)
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
