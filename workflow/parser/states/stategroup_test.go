package states

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseStateGroup(t *testing.T) {

	var testcases = []struct {
		definition string
	}{
		{
			definition: `
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
			}
			`,
		}, {
			definition: `
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
			  }`,
		},
	}

	for idx, testcase := range testcases {
		fmt.Println("idx -- ", idx)
		pl, err := NewStateGroupFromString(testcase.definition, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(pl)
		bone := pl.GetBone()
		bonebyte, err := json.Marshal(bone)
		fmt.Println(string(bonebyte))
	}

}
