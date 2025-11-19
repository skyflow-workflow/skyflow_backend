package stepfunction

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/config"
)

func TestParserStatemachineHeaderFailed(t *testing.T) {
	var testcases = []struct {
		name      string
		defintion string
	}{
		{
			name:      "version int",
			defintion: `{"version":1.0, "type":"stepfunction"}`,
		},
		{
			name:      "bad type",
			defintion: `{"version":false, "type":"stepfunction"}`,
		},
	}
	decoder := NewStepFunctionDecoder(&config.StandardExecutorConfig)
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			header, err := decoder.DecodeStateMachineHeaderDefinition(tt.defintion)
			assert.NotEqual(t, err, nil)
			if err == nil {
				t.Log(header)
			}
		})
	}
}
