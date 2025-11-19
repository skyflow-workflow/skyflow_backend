package states

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/expression/stepfunction"
)

func TestValidateParsedChoice(t *testing.T) {
	var testcases = []struct {
		name        string
		choice      *ChoiceState
		input       map[string]any
		expectNext  string
		expectError error
	}{
		{
			name: "choice with string equals",
			choice: &ChoiceState{
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
			choice: &ChoiceState{
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
			choice: &ChoiceState{
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

func TestParserChoiceState(t *testing.T) {

	var testcases = []struct {
		template  string
		reason    string
		wantError bool
	}{
		{
			template: `	{
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
			wantError: false,
		},
		{
			template: `	{
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
				"Default": 10
			  }`,
			wantError: true,
		},
		{
			template: `	{
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
			wantError: false,
		},
		{
			template: `	{
				"Type" : "Choice",
				"Choices": [
				],
				"Default": ""
			  }`,
			wantError: true,
		},
		{
			template: `	{
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
			wantError: false,
		},
	}
	for _, tt := range testcases {
		_, err := NewChoiceStateFromString(tt.template)
		assert.Equal(t, tt.wantError, err != nil)

	}
}

func TestRunChoice(t *testing.T) {

	var testcases = []struct {
		name      string
		defintion string
		input     string
		wantBone  StateBone
		wantNext  string
		wantError bool
	}{
		{
			name: "choice with one default",
			defintion: `{
				"Type" : "Choice",
				"Comment": "2 choice with one default",
				"Choices": [
				  {
					"Variable": "$.foo",
					"NumericGreaterThan": 3,
					"Next": "FirstMatchState"
				  },
				  {
					"Variable": "$.boo",
					"NumericEquals": 2,
					"Next": "SecondMatchState"
				  }
				],
				"Default": "DefaultState"
			  }`,
			input: `{
				"foo": 2.0,
				"boo": 2
			  }`,
			wantBone: StateBone{
				BaseBone: BaseBone{
					Type:    "Choice",
					Name:    "",
					Next:    []string{"FirstMatchState", "SecondMatchState", "DefaultState"},
					Comment: "2 choice with one default",
					End:     false,
				},
			},
			wantNext:  "SecondMatchState",
			wantError: true,
		},
		{
			name: "choice without default and two choice",
			defintion: `{
					"Type": "Choice",
					"Choices": [
						{
							"Variable": "$.snapflags",
							"VariableExist": true,
							"Next": "MTNCGetBoxID"
						},
						{
							"Variable": "$.snapflags",
							"NumericEquals": 0,
							"Next": "MTNCEnableDepotNoBox"
						}
					],
					"Comment": "根据需要决定是否检查快照"
				}`,
			input: `{
					"foo":      2.0,
					"boo":      2,
					"snapflags": false
				}`,
			wantBone: StateBone{
				BaseBone: BaseBone{
					Type:    "Choice",
					Name:    "",
					Next:    []string{"MTNCGetBoxID", "MTNCEnableDepotNoBox"},
					Comment: "根据需要决定是否检查快照",
					End:     false,
				},
			},
			wantNext:  "MTNCGetBoxID",
			wantError: true,
		},
		{
			name: "2 choice with one default",
			defintion: `{
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
				  },
				  {
					"Variable": "$.foo",
					"IsNull": false,
					"Next": "WaitState2"
				  },
				  {
					"Next": "Pass2",
					"Variable": "$.foo",
					"IsExist": false
				  }
				],
				"Default": "DefaultState",
				"InputPath": "$.boo",
				"Type": "Choice"
			  }`,
			wantBone: StateBone{
				BaseBone: BaseBone{
					Type:    "Choice",
					Name:    "",
					Next:    []string{"FirstMatchState", "SecondMatchState", "WaitState2", "Pass2", "DefaultState"},
					Comment: "",
					End:     false,
				},
			},
			input:     `{"boo":{"foo":2,"foo2":3}}`,
			wantNext:  "SecondMatchState",
			wantError: false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			cc, err := NewChoiceStateFromString(tt.defintion)
			assert.Equal(t, err == nil, true)
			bone := cc.GetBone()
			assert.Equal(t, tt.wantBone, bone)
			inputmap, err := StringToMap(tt.input)
			assert.Equal(t, err == nil, true)
			next, err := cc.GetNextState(inputmap)
			assert.Equal(t, err == nil, true)
			assert.Equal(t, tt.wantNext, next.Name)
		})

	}
}
