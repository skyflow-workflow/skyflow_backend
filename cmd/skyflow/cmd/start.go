package cmd

import (
	"fmt"

	"github.com/skyflow-workflow/skyflow_backbend/config"
	"github.com/skyflow-workflow/skyflow_backbend/server/apiserver"
	"github.com/spf13/cobra"
	tconfig "trpc.group/trpc-go/trpc-go/config"
)

var (
	trpc_conf string = "./trpc_go.yaml" // Configuration file for the TRPC server
)

func init() {
	startCmd.Flags().StringVarP(&trpc_conf, "trpc_conf", "c", "./trpc_go.yaml", " trpc configuration file, default is ./trpc_go.yaml")

}

var startCmd = &cobra.Command{
	Use:   "start [flags] [args]",
	Short: "Start the skyflow workflow engine",
	Long:  `Start the skyflow workflow engine, which allows you to run workflows defined in JSON format.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(trpc_conf)
	},
}

func StartAPIServer() {
	// Initialize the API server
	apiServer := apiserver.NewApiServer(trpc_conf)
	// Start the API server
	apiServer.Start()
}

// func StartDispatcher() {
// 	// Initialize the dispatcher
// 	dispatcher := dispatcher.NewDispatcher(trpc_conf)

// 	// Start the dispatcher
// 	if err := dispatcher.Start(); err != nil {
// 		fmt.Println("Failed to start dispatcher:", err)
// 	}
// }

func LoadService(customConfigFilePath string) error {

	customCfg, err := tconfig.Load(customConfigFilePath, tconfig.WithCodec("yaml"))
	if err != nil {
		return fmt.Errorf("error loading custom configuration file: %w", err)
	}
	var customConfig = config.SkyflowConfig{}
	err = customCfg.Unmarshal(&customConfig)
	if err != nil {
		return fmt.Errorf("error unmarshalling custom configuration: %w", err)
	}
	return nil

}
