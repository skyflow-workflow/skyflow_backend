package states

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/expression/stepfunction"
)

func TestChoice(t *testing.T) {
	var testcases = []struct {
		name        string
		choice      *Choice
		input       map[string]any
		expectNext  string
		expectError error
	}{
		{
			name: "choice with string equals",
			choice: &Choice{
				BaseState: &BaseState{
					Name: "choice",
					Type: "Choice",
				},
				ChoiceBody: &ChoiceBody{
					Choices: []ChoiceBranch{
						{
							Condition: &stepfunction.EvaluateCondition{
								EvaluateUnit: &stepfunction.EvaluateUnit{
									Variable: "$.foo",
									Operator: "StringEquals",
									Operand:  "bar",
								},
							},
							Next: "branch1",
						},
					},
					Default: "default_next",
				},
			},
			input: map[string]any{
				"foo": "bar",
			},
			expectNext:  "branch1",
			expectError: nil,
		},
		{
			name: "choice with string notequals",
			choice: &Choice{
				BaseState: &BaseState{
					Name: "choice",
					Type: "Choice",
				},
				ChoiceBody: &ChoiceBody{
					Choices: []ChoiceBranch{
						{
							Condition: &stepfunction.EvaluateCondition{
								EvaluateUnit: &stepfunction.EvaluateUnit{
									Variable: "$.foo",
									Operator: "StringEquals",
									Operand:  "bar",
								},
							},
							Next: "branch1",
						},
					},
					Default: "default_next",
				},
			},
			input: map[string]any{
				"foo": "notbar",
			},
			expectNext:  "default_next",
			expectError: nil,
		},
		{
			name: "choice with empty choices",
			choice: &Choice{
				BaseState: &BaseState{
					Name: "choice",
					Type: "Choice",
				},
				ChoiceBody: &ChoiceBody{
					Choices: []ChoiceBranch{},
					Default: "default_next",
				},
			},
			input: map[string]any{
				"foo": "notbar",
			},
			expectNext:  "default_next",
			expectError: nil,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			next, err := testcase.choice.GetNextState(testcase.input)
			if err != nil {
				assert.Equal(t, testcase.expectError, nil)
				assert.Equal(t, testcase.expectError.Error(), err.Error())
				return
			}
			assert.Equal(t, testcase.expectNext, next.Name)
		})
	}
}
