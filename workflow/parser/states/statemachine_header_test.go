package states

import (
	"fmt"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestParserHeader(t *testing.T) {

	var testcases = []struct {
		template  string
		wantError bool
	}{
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": 10,
				"ExecutionInfoPath": "$.execution_info"
				}`,
			wantError: false,
		},
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": -1,
				"ExecutionInfoPath": "$.execution_info"
				}`,
			wantError: true,
		},
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": 1,
				"ExecutionInfoPath": "$$.execution_info"
				}`,
			wantError: true,
		},
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": 1,
				"ExecutionInfoPath": "$.execution_info",
				"DefaultInput" :10
				}`,
			wantError: true,
		},
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": 1,
				"ExecutionInfoPath": "$.execution_info",
				"DefaultInput" :{
						"abc": 10
					}
				}`,
			wantError: false,
		},
	}

	for idx, tt := range testcases {
		fmt.Println("process index : ", idx)
		state, err := NewStateMachineHeaderFromString(tt.template)
		fmt.Println(err)
		fmt.Println(state)
		// asset
		assert.Equal(t, err != nil, tt.wantError)

	}
}

func TestParserHeaderInput(t *testing.T) {

	var testcases = []struct {
		template  string
		input     string
		wantError bool
	}{
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": 10,
				"ExecutionInfoPath": "$.execution_info"
				}`,
			input: `{
				"k1":"v1"
			}`,
			wantError: false,
		},
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": 1,
				"ExecutionInfoPath": "$.execution_info",
				"DefaultInput" :{
						"abc": 10
					}
				}`,
			input: `{
				"k1":"v1"
			}`,
			wantError: false,
		},
		{
			template: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": 1,
				"ExecutionInfoPath": "$.execution_info",
				"DefaultInput" :{
						"defaultkey": 10
					}
				}`,
			input: `{
				"k1":"v1" ,
				"defaultkey": {
					"jobid":"xxxxmockjobid"
				}
			}`,
			wantError: false,
		},
	}
	var mockexecution = map[string]interface{}{
		"url":  "mockurl",
		"name": "mockname",
	}
	for idx, tt := range testcases {
		fmt.Println("process index : ", idx)
		state, err := NewStateMachineHeaderFromString(tt.template)
		fmt.Println(err)
		fmt.Println(state)
		// asset
		assert.Equal(t, err == nil, true)
		input, err := StringToMap(tt.input)
		assert.Equal(t, err == nil, true)
		innerinput, err := state.GetInput(input, mockexecution)
		assert.Equal(t, err != nil, tt.wantError)
		fmt.Println(innerinput)
	}
}
