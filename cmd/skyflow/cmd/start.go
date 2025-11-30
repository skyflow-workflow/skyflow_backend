package cmd

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/goodaye/wire"
	"github.com/skyflow-workflow/skyflow_backend/cmd/skyflow/kratos"
	"github.com/skyflow-workflow/skyflow_backend/config"
	"github.com/skyflow-workflow/skyflow_backend/server/dispatcher"
	"github.com/skyflow-workflow/skyflow_backend/workflow"
	"github.com/spf13/cobra"
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
	frameworkSever, err := NewFrameworkAPIServer(frame_conf, wfSvc)
	if err != nil {
		slog.Error("Error loading framework server", "error", err)
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
				if err := frameworkSever.Start(); err != nil {
					slog.Error("Error starting API server", "error", err)
					return
				}
			}()
		case "dispatcher":
			wg.Add(1)
			go func() {
				defer wg.Done()
				slog.Info("Starting Dispatcher...")
				StartDispatcher(wfSvc, sfConfig.Dispatcher)
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

func StartDispatcher(wfSvc workflow.WorkflowService, conf *config.DispatcherConfig) {
	// Initialize the dispatcher
	dispatcher, err := dispatcher.NewDispatcher(wfSvc, conf)
	if err != nil {
		slog.Error("Failed to initialize dispatcher:", "error", err)
		panic(err)
	}
	// Start the dispatcher
	if err = dispatcher.Start(); err != nil {
		slog.Error("Failed to start dispatcher", "error", err)
		panic(err)
	}
}

func NewFrameworkAPIServer(framework_conf string, wfSvc workflow.WorkflowService) (FrameWorkAPIServer, error) {

	// if framework  is  "trpc"
	// server, err := trpc.NewServer(framework_conf, wfSvc)
	// if err != nil {
	// 	return nil, err
	// }

	// if framework_conf == "kratos"
	server, err := kratos.NewServer(framework_conf, wfSvc)
	if err != nil {
		return nil, err
	}

	return server, nil
}
