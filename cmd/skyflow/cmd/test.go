package cmd

import (
	"log/slog"
	"os"

	"github.com/skyflow-workflow/skyflow_backbend/config"
	"github.com/spf13/cobra"
	"trpc.group/trpc-go/trpc-go"
	tconfig "trpc.group/trpc-go/trpc-go/config"
)

var customConfigFilePath = "./skyflow.yaml"

func init() {
	testCmd.Flags().StringVarP(&trpc_conf, "trpc_conf", "c", "./trpc_go.yaml", " trpc configuration file, default is ./trpc_go.yaml")

}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test config validation for start skyflow workflow engine",
	Long:  `Test config validation for start skyflow workflow engine.`,
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		err = CheckConfig()
		if err != nil {
			slog.Error("Error checking configuration", "error", err)
			os.Exit(1)
		}
		slog.Info("Configuration check completed successfully")

	},
}

func CheckConfig() error {
	if trpc_conf != "" {
		trpc.ServerConfigPath = trpc_conf // Set the TRPC server configuration path
	}
	_, err := trpc.LoadConfig(trpc_conf)
	if err != nil {
		slog.Error("Error loading TRPC server configuration", "error", err, "configPath", trpc_conf)
		return err
	}
	_ = trpc.NewServer()
	c, err := tconfig.Load(customConfigFilePath, tconfig.WithCodec("yaml"))
	if err != nil {
		slog.Error("Error loading custom configuration file", "error", err, "configPath", customConfigFilePath)
		return err
	}
	var customConfig = config.SkyflowConfig{}
	err = c.Unmarshal(&customConfig)
	if err != nil {
		slog.Error("Error unmarshalling custom configuration", "error", err, "configPath", customConfigFilePath)
		return err
	}
	return nil
}
