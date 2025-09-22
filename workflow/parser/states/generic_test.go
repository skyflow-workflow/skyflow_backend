package states

import (
	"fmt"
	"testing"
)

func TestUseGeneric(t *testing.T) {

	taskjson := `{
		"Type": "Task",
		"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME",
		"InputPath" : "$.mydata",
		"Parameters" : {
			"abc.$" : "$.foo",
			"123": [{
				"var1.$" : "$.boo",
				"var2.$" : "$.foo"
			}
			]
		},
		"MaxExecuteTimes":10,
		"Retry":[
			{
			  "ErrorEquals": [ "ErrorA", "ErrorB" ],
			  "IntervalSeconds": 1,
			  "BackoffRate": 2.0,
			  "MaxAttempts": 2
			},
			{
			  "ErrorEquals": [ "ErrorC" ],
			  "IntervalSeconds": 5
			}
		  ],
		  "Catch": [
			{
			  "ErrorEquals": [ "States.ALL" ],
			  "Next": "Z"
			}
		  ],
		"Next": "ChoiceState"
	  }`

	input := map[string]interface{}{
		"mydata": map[string]interface{}{
			"foo": 2.0,
			"boo": 2,
		},
	}
	// var taskmap map[string]interface{}

	// err := json.Unmarshal([]byte(taskjson), &taskmap)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("%v", taskmap)
	// cc := state.NewTaskState("abc")
	// err = cc.Inititalize(taskmap)

	cc, err := NewGenericStateFromString[TaskState](taskjson, NewTaskStateFromMap)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", cc)
	fmt.Println(cc.MaxExecuteTimes)
	fmt.Println(cc.Parameters)
	next, err := cc.GetInput(input)
	fmt.Println(next, err)
	fmt.Println(cc.GetBone())

}
