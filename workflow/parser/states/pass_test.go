package states

import (
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestPassGetResult(t *testing.T) {
	testcases := []struct {
		name          string
		state         *PassState
		input         any
		expected      any
		expectedError error
	}{
		{
			name: "test simple pass",
			state: &PassState{
				BaseState: &BaseState{
					Name: "test",
				},
				PassBody: &PassBody{
					Result: map[string]any{
						"result_key1": "result_value1",
						"result_key2": "result_value2",
					},
				},
			},
			input: map[string]any{
				"input_key1": "input_value1",
				"input_key2": "input_value2",
			},
			expected: map[string]any{
				"result_key1": "result_value1",
				"result_key2": "result_value2",
			},
		},
		{
			name: "test pass with parameters",
			state: &PassState{
				BaseState: &BaseState{
					Name: "test",
				},
				PassBody: &PassBody{
					Result: map[string]any{
						"result_key1.$": "$.input_key1",
						"result_key2.$": "$.input_key2",
					},
				},
			},
			input: map[string]any{
				"input_key1": "input_value1",
				"input_key2": "input_value2",
			},
			expected: map[string]any{
				"result_key1": "input_value1",
				"result_key2": "input_value2",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.state.GetResult(tc.input)
			if err != nil {
				assert.Equal(t, tc.expectedError != nil, true)
				assert.Equal(t, err.Error(), tc.expectedError.Error())
				return
			}
			assert.Equal(t, result, tc.expected)
		})
	}
}

func TestPass_GetDefinition(t *testing.T) {

	rawdata := `
	{
		"Type": "Pass",
		"Result": {
		  "result_key1.$": "$.input_key1",
		  "result_key2.$": "$.input_key2"
		},
		"ResultPath": "$.result",
		"Next": "NextState"
	  }
	`
	state, err := NewPassStateFromString(rawdata)
	if err != nil {
		t.Fatal(err)
	}
	def, err := state.GetDefinition()
	assert.Equal(t, err, nil)
	// GetDefinition 会过滤 Deny 字段（Retry、Catch），且字段顺序不固定，故只校验必要内容与不包含 Deny 字段
	assert.Equal(t, strings.Contains(def, "Type"), true)
	assert.Equal(t, strings.Contains(def, "Pass"), true)
	assert.Equal(t, strings.Contains(def, "ResultPath"), true)
	assert.Equal(t, strings.Contains(def, "Next"), true)
	assert.Equal(t, strings.Contains(def, "Result"), true)
	assert.Equal(t, strings.Contains(def, "Retry"), false)
	assert.Equal(t, strings.Contains(def, "Catch"), false)
}
