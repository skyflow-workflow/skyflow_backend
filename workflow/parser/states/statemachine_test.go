package states

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestParseMapStateMachine(t *testing.T) {

	var testcases = []struct {
		template  string
		wantError bool
	}{
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"Type": "stepfunction",
				"TimeoutSeconds": 1243,
				"States": {
				  "FirstState": {
					"Type": "Task",
					"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME",
					"Next": "ChoiceState"
				  },
				  "ChoiceState": {
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
				  },

				  "FirstMatchState": {
					"Type" : "Task",
					"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:OnFirstMatch",
					"Next": "NextState"
				  },

				  "SecondMatchState": {
					"Type" : "Task",
					"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:OnSecondMatch",
					"Next": "NextState"
				  },

				  "DefaultState": {
					"Type": "Fail",
					"Error": "DefaultStateError",
					"Cause": "No Matches!"
				  },

				  "NextState": {
					"Type": "Task",
					"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME",
					"End": true
				  }
				}
			  }
		`,
			wantError: false,
		},
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": 1243,
				"States": {
				  "FirstState": {
					"Type": "Task",
					"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME",
					"End": true
				  }
				}
			}
			`,
			wantError: false,
		},
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": 1,
				"TimeoutSeconds": 1243,
				"States": {
				  "FirstState": {
					"Type": "Task",
					"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME",
					"End": true
				  }
				}
			}
			`,
			wantError: true,
		},
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": -1243,
				"States": {
				  "FirstState": {
					"Type": "Task",
					"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME",
					"End": true
				  }
				}
			}
			`,
			wantError: true,
		},
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": 0,
				"_Bone":"abc",
				"States": {
				  "FirstState": {
					"Type": "Task",
					"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME",
					"End": true
				  }
				}
			}
			`,
			wantError: false,
		},
		{
			template: `
			{
				"Comment": "Parallel Example.",
				"StartAt": "LookupCustomerInfo",
				"Type": "stepfunction",
				"States": {
				  "LookupCustomerInfo": {
					"Type": "Parallel",
					"End": true,
					"Branches": [
					  {
					   "StartAt": "LookupAddress",
					   "States": {
						 "LookupAddress": {
						   "Type": "Task",
						   "Resource":
							 "arn:aws-cn:lambda:us-east-1:123456789012:function:AddressFinder",
						   "End": true
						 }
					   }
					 },
					 {
					   "StartAt": "LookupPhone",
					   "States": {
						 "LookupPhone": {
						   "Type": "Task",
						   "Resource":
							 "arn:aws-cn:lambda:us-east-1:123456789012:function:PhoneFinder",
						   "End": true
						 }
					   }
					 }
					]
				  }
				}
			  }
			`,
			wantError: false,
		},
	}

	for idx, tt := range testcases {
		fmt.Printf("process index --- %d\n", idx)
		var data map[string]interface{}
		err := json.Unmarshal([]byte(tt.template), &data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		sm, err := NewStateMachineFromMap(data)
		fmt.Println("err :", err)
		assert.Equal(t, err != nil, tt.wantError)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// fmt.Println(sm)
		// fmt.Printf("%+v", sm)
		bone := sm.GetBone()
		fmt.Printf(" bone  index %d : \n", idx)
		fmt.Println(bone)
		groupstate, err := sm.GetGroupStates()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf(" group state  index %d : \n", idx)
		fmt.Println(ToString(groupstate))
	}

}

func TestStateMachineBone(t *testing.T) {
	var testcases = []struct {
		template  string
		wantError bool
	}{
		{
			template: `
			{
				"StartAt":"P1",
				"States":{
					"P1":{
						"Type":"Pass",
						"End":true
					}
				}
			}
			`,
		},
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"States": {
				  "FirstState": {
					"Type": "Pass",
					"Next": "ChoiceState"
				  },
				  "ChoiceState": {
					"InputPath": "$.boo",
					"Type": "Choice",
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
				  },
				  "FirstMatchState": {
					"Type": "Pass",
					"Next": "WaitState"
				  },
				  "SecondMatchState": {
					"Type": "Pass",
					"Next": "NextState"
				  },
				  "DefaultState": {
					"Type": "Fail",
					"Error": "DefaultStateError",
					"Cause": "No Matches!"
				  },
				  "NextState": {
					"Type": "Succeed"
				  },
				  "WaitState": {
					"Type": "Wait",
					"InputPath": "$.foo.x",
					"Seconds": 10,
					"Next": "NextState"
				  }
				}
			  }`,
		},
	}

	for _, tt := range testcases {
		sm, err := NewStateMachineFromString(tt.template)
		assert.Equal(t, err, nil)
		bone := sm.GetBone()
		fmt.Println(bone)
	}
}
