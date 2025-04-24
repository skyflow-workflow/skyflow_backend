package states

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestPassGetResult(t *testing.T) {
	testcases := []struct {
		name          string
		state         *Pass
		input         any
		expected      any
		expectedError error
	}{
		{
			name: "test simple pass",
			state: &Pass{
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
			state: &Pass{
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
