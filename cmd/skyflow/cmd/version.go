package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// Version of the skyflow workflow engine
var version string = "0.1.0"

var VersionOutput = struct {
	Version string `json:"version"` // Version of the skyflow workflow engine
}{
	Version: version,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of skyflow workflow engine",
	Long:  `Print the version number of skyflow workflow engine.`,
	Run: func(cmd *cobra.Command, args []string) {

		var output string
		if format == outputformats.JSON {
			outputbyte, err := json.Marshal(VersionOutput)
			if err != nil {
				slog.Error("Failed to marshal version output", "error", err)
				os.Exit(1)
			}
			output = string(outputbyte)
		} else {
			output = fmt.Sprintf("version: %s", version)
		}
		fmt.Println(output)
	},
}
