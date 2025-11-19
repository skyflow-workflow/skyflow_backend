package cmd

import (
	"log/slog"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestCheckConfig(t *testing.T) {
	// This function is a placeholder for the actual test implementation.
	// It should contain the logic to test the CheckConfig function.
	// For now, we will just call the function to ensure it runs without errors.
	var err error
	wd, err := os.Getwd()
	assert.Equal(t, err, nil)
	slog.Info("Current working directory:", "wd", wd)
	frame_conf = "./mock/trpc_go.yaml"
	skyflow_conf = "./mock/skyflow.yaml"
	err = CheckConfig()
	assert.Equal(t, err, nil)

}
