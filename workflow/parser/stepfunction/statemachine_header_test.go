package stepfunction

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/decoder"
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
	decoder := NewStepfuncionDecoder(&decoder.StandardParserConfig, &decoder.DefaultQuota)
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			header, err := decoder.DecodeStateMachineHeaderDefintion(tt.defintion)
			assert.NotEqual(t, err, nil)
			if err == nil {
				t.Log(header)
			}
		})
	}
}
