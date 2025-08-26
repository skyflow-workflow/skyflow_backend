package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var syncSchemaCmd = &cobra.Command{
	Use:   "syncschema [flags] [args]",
	Short: "Sync workflow schema to the Database/MessageQueue",
	Long:  `Sync workflow schema to the Database/MessageQueue`,
	Run:   SyncSchemaCommand,
}

func SyncSchemaCommand(cmd *cobra.Command, args []string) {

	sfConfig, err := LoadConfig(skyflow_conf)
	if err != nil {
		slog.Error("Error loading skyflow configuration", "error", err)
		return
	}
	slog.Info("Skyflow configuration loaded successfully")
	wfSvc, err := LoadService(sfConfig)
	if err != nil {
		slog.Error("Error loading services", "error", err)
		return
	}
	slog.Info("Services loaded successfully")

	err = wfSvc.SyncSchema()
	if err != nil {
		slog.Error("Error syncing workflow schema", "error", err)
		return
	}
	slog.Info("Workflow schema synced successfully")
}
