package stepfunction

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/expression/stepfunction"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"

	"gopkg.in/go-playground/assert.v1"
)

func TestDecodeChoiceState(t *testing.T) {

	var testcases = []struct {
		name       string
		definition string
		expected   *states.Choice
		wantError  error
	}{
		{
			name: "choice with two numeric equals",
			definition: `	{
				"Type" : "Choice",
				"Choices": [
				  {
					"Variable": "$.foo",
					"NumericEquals": 1,
					"Next": "FirstMatchState"
				  },
				  {
					"Variable": "$.foo",
					"NumericEquals": 2,
					"Next": "SecondMatchState"
				  }
				],
				"Default": "DefaultState"
			  }`,
			expected: &states.Choice{
				BaseState: &states.BaseState{
					Type:            "Choice",
					InputPath:       "$",
					OutputPath:      "$",
					ResultPath:      "$",
					MaxExecuteTimes: 1000,
				},
				ChoiceBody: &states.ChoiceBody{
					Choices: []states.ChoiceBranch{
						{
							Condition: &stepfunction.EvaluateCondition{
								EvaluateUnit: &stepfunction.EvaluateUnit{
									Variable: "$.foo",
									Operator: "NumericEquals",
									Operand:  float64(1),
								},
							},
							Next: "FirstMatchState",
						},
						{
							Condition: &stepfunction.EvaluateCondition{
								EvaluateUnit: &stepfunction.EvaluateUnit{
									Variable: "$.foo",
									Operator: "NumericEquals",
									Operand:  float64(2),
								},
							},
							Next: "SecondMatchState",
						},
					},
					Default: "DefaultState",
				},
			},
			wantError: nil,
		},
		{
			name: "choice with empty choices",
			definition: `	{
				"Type" : "Choice",
				"Choices": [],
				"Default": "DefaultState"
			  }`,
			expected: &states.Choice{
				BaseState: &states.BaseState{
					Type:            "Choice",
					InputPath:       "$",
					OutputPath:      "$",
					ResultPath:      "$",
					MaxExecuteTimes: 1000,
				},
				ChoiceBody: &states.ChoiceBody{
					Choices: []states.ChoiceBranch{},
					Default: "DefaultState",
				},
			},
			wantError: nil,
		},
		{
			name: "choice with empty default",
			definition: `	{
				"Type" : "Choice",
				"Choices": [
				  {
					"Variable": "$.foo",
					"NumericEquals": 1,
					"Next": "FirstMatchState"
				  },
				  {
					"Variable": "$.foo",
					"NumericEquals": 2,
					"Next": "SecondMatchState"
				  }
				],
				"Default": ""
			  }`,
			expected: &states.Choice{
				BaseState: &states.BaseState{
					Type:            "Choice",
					InputPath:       "$",
					OutputPath:      "$",
					ResultPath:      "$",
					MaxExecuteTimes: 1000,
				},
				ChoiceBody: &states.ChoiceBody{
					Choices: []states.ChoiceBranch{
						{
							Condition: &stepfunction.EvaluateCondition{
								EvaluateUnit: &stepfunction.EvaluateUnit{
									Variable: "$.foo",
									Operator: "NumericEquals",
									Operand:  float64(1),
								},
							},
							Next: "FirstMatchState",
						},
						{
							Condition: &stepfunction.EvaluateCondition{
								EvaluateUnit: &stepfunction.EvaluateUnit{
									Variable: "$.foo",
									Operator: "NumericEquals",
									Operand:  float64(2),
								},
							},
							Next: "SecondMatchState",
						},
					},
					Default: "",
				},
			},
			wantError: nil,
		},
		{
			name: "choice with empty choices and default",
			definition: `	{
				"Type" : "Choice",
				"Choices": [
				],
				"Default": ""
			  }`,
			expected: &states.Choice{
				BaseState: &states.BaseState{
					Type:            "Choice",
					InputPath:       "$",
					OutputPath:      "$",
					ResultPath:      "$",
					MaxExecuteTimes: 1000,
				},
				ChoiceBody: &states.ChoiceBody{
					Choices: []states.ChoiceBranch{},
					Default: "",
				},
			},
			wantError: nil,
		},
		{
			name: "choice with multiple conditions",
			definition: `	{
				"Type": "Choice",
				"Choices": [
				  {
					  "Not": {
						"Variable": "$.type",
						"StringEquals": "Private"
					  },
					  "Next": "Public"
				  },
				  {
					"Variable": "$.value",
					"NumericEquals": 0,
					"Next": "ValueIsZero"
				  },
				  {
					"And": [
					  {
						"Variable": "$.value",
						"NumericGreaterThanEquals": 20
					  },
					  {
						"Variable": "$.value",
						"NumericLessThan": 30
					  }
					],
					"Next": "ValueInTwenties"
				  }
				],
				"Default": "DefaultState"
			  }`,
			expected: &states.Choice{
				BaseState: &states.BaseState{
					Type:            "Choice",
					InputPath:       "$",
					OutputPath:      "$",
					ResultPath:      "$",
					MaxExecuteTimes: 1000,
				},
				ChoiceBody: &states.ChoiceBody{
					Choices: []states.ChoiceBranch{
						{
							Condition: &stepfunction.EvaluateCondition{
								Not: &stepfunction.EvaluateCondition{
									EvaluateUnit: &stepfunction.EvaluateUnit{
										Variable: "$.type",
										Operator: "StringEquals",
										Operand:  "Private",
									},
								},
							},
							Next: "Public",
						},
						{
							Condition: &stepfunction.EvaluateCondition{
								EvaluateUnit: &stepfunction.EvaluateUnit{
									Variable: "$.value",
									Operator: "NumericEquals",
									Operand:  float64(0),
								},
							},
							Next: "ValueIsZero",
						},
						{
							Condition: &stepfunction.EvaluateCondition{
								And: []stepfunction.EvaluateCondition{
									{
										EvaluateUnit: &stepfunction.EvaluateUnit{
											Variable: "$.value",
											Operator: "NumericGreaterThanEquals",
											Operand:  float64(20),
										},
									},
									{
										EvaluateUnit: &stepfunction.EvaluateUnit{
											Variable: "$.value",
											Operator: "NumericLessThan",
											Operand:  float64(30),
										},
									},
								},
							},
							Next: "ValueInTwenties",
						},
					},
					Default: "DefaultState",
				},
			},
			wantError: nil,
		},
	}
	decoder := NewStepfuncionDecoder(nil, nil)

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			err = tt.expected.BaseState.Init()
			assert.Equal(t, err, nil)

			state, err := decoder.DecodeStateDefintion(tt.definition)
			assert.Equal(t, err, nil)
			opts := cmp.Options{cmpopts.IgnoreUnexported(states.BaseState{}), cmpopts.IgnoreUnexported(states.TaskBody{})}
			diff := cmp.Diff(state, tt.expected, opts...)
			assert.Equal(t, "", diff)

		})
	}
}
