package cmd

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestCheckConfig(t *testing.T) {
	// This function is a placeholder for the actual test implementation.
	// It should contain the logic to test the CheckConfig function.
	// For now, we will just call the function to ensure it runs without errors.
	trpc_conf = "./mock/trpc_go_dev.yaml"
	customConfigFilePath = "./mock/custom.yaml"
	err := CheckConfig()
	assert.Equal(t, err, nil)

}
