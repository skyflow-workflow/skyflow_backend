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

// NewParallelStateFromString NewParallelStateFromString
func NewParallelStateFromString(definition string, depth int) (*Parallel, error) {

	state, err := NewGenericGroupStateFromString(definition, depth, NewParallelStateFromMap)
	return state, err
}

// NewParallelStateFromMap NewParallelStateFromMap
func NewParallelStateFromMap(data map[string]interface{}, depth int) (*Parallel, error) {

	var err error
	state := &Parallel{
		_Branches: []*StateGroup{},
		_Depth:    depth,
	}

	// basestate
	bs, err := NewBaseStateFromMap(data)
	if err != nil {
		return state, err
	}
	state.BaseState = bs

	err = state.InitByMap(data)
	return state, err
}

// InitByMap Inititalize TaskState Content
func (s *Parallel) InitByMap(data map[string]interface{}) error {
	var err error

	// 初始化自身
	err = MapStructDecode(data, s)
	if err != nil {
		return err
	}

	err = myvalidate.Struct(s)
	if err != nil {
		return err
	}
	err = s.Init()
	if err != nil {
		return err
	}
	return err
}
