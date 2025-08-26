package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	format              = outputformats.Text // Output format for the command, e.g., json, yaml
	trpc_conf    string = "./trpc_go.yaml"   // Configuration file for the TRPC framework
	skyflow_conf string = "./skyflow.yaml"   // Configuration file for the skyflow server
)

var outputformats = struct {
	Text string
	JSON string
}{
	Text: "text",
	JSON: "json",
}

func init() {
	cobra.OnInitialize(initCobra)
	rootCmd.CompletionOptions.DisableDefaultCmd = true // Disable default completion command
	rootCmdAddFlag()
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(syncSchemaCmd)

}

func initCobra() {

	if format == outputformats.JSON {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})))
	}
}

func rootCmdAddFlag() {
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", outputformats.Text,
		"Output format for the command. Supported formats: text, json. Default is 'text'.")
	rootCmd.PersistentFlags().SortFlags = false // Disable sorting of flags
	rootCmd.Flags().StringVarP(&trpc_conf, "trpc_config", "c", "./trpc_go.yaml", " trpc configuration file, default is ./trpc_go.yaml")
	rootCmd.Flags().StringVarP(&skyflow_conf, "skyflow_config", "", "./skyflow.yaml", " skyflow configuration file, default is ./skyflow.yaml")

}

var rootCmd = &cobra.Command{
	Use:   "skyflow",
	Short: "skyflow is a json based workflow engine",
	Long: `A json-based workflow engine that allows you to define workflows in a declarative manner.
	workflow DSL defined in AWS Step Functions, but with more flexibility and extensibility.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to skyflow workflow engine!")
		fmt.Println("Use 'skyflow --help' to see available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Failed to execute command", "error", err)
		os.Exit(1)
	}
}
