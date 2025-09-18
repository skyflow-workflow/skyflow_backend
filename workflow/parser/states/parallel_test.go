package states

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseParallelState(t *testing.T) {

	template := ` {
				"Type": "Parallel",
				"Comment": "Lookup Address and Phone",
				"End": true,
				"Branches": [
				  {
				   "Comment": "Lookup Address",
				   "InputPaht": "$.address",
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
					"Comment": "Lookup Address",
					"InputPaht": "$.phone",
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

	`
	pl, err := NewParallelStateFromString(template, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(pl)
	bone := pl.GetBone()
	bonestr, _ := json.Marshal(bone)
	fmt.Println(string(bonestr))

}
