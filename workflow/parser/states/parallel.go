package states

import (
	"fmt"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

/*
 * @Author: mumangtao@gmail.com
 * @Date: 2019-12-05 20:38:26
 * @Last Modified by: mumangtao@gmail.com
 * @Last Modified time: 2019-12-30 11:37:00
 */
// Parallel State

// {
// 	"Comment": "Parallel Example.",
// 	"StartAt": "LookupCustomerInfo",
// 	"States": {
// 	  "LookupCustomerInfo": {
// 		"Type": "Parallel",
// 		"End": true,
// 		"Branches": [
// 		  {
// 		   "StartAt": "LookupAddress",
// 		   "States": {
// 			 "LookupAddress": {
// 			   "Type": "Task",
// 			   "Resource":
// 				 "arn:aws-cn:lambda:us-east-1:123456789012:function:AddressFinder",
// 			   "End": true
// 			 }
// 		   }
// 		 },
// 		 {
// 		   "StartAt": "LookupPhone",
// 		   "States": {
// 			 "LookupPhone": {
// 			   "Type": "Task",
// 			   "Resource":
// 				 "arn:aws-cn:lambda:us-east-1:123456789012:function:PhoneFinder",
// 			   "End": true
// 			 }
// 		   }
// 		 }
// 		]
// 	  }
// 	}
//   }

// ParallelState 并行状态
type ParallelState struct {
	*BaseState    `json:",inline"`
	*ParallelBody `json:",inline"`
	_Depth        int
}

type ParallelBody struct {
	Branches  []map[string]interface{} `mapstructure:"Branches" validate:"required,gt=0"`
	_Branches []*StateMachineBody
}

// NewParallelStateFromString NewParallelStateFromString
func NewParallelStateFromString(definition string, depth int) (*ParallelState, error) {

	// state, err = NewGenericStateFromString(definition, NewTaskStateFromMap)
	mapdata, err := StringToMap(definition)
	if err != nil {
		return nil, err
	}

	state, err := NewParallelStateFromMap(mapdata, depth)
	if err != nil {
		return nil, err
	}
	return state, nil

}

// NewParallelStateFromMap NewParallelStateFromMap
func NewParallelStateFromMap(data map[string]interface{}, depth int) (*ParallelState, error) {

	// basestate
	bs, err := NewBaseStateFromMap(data)
	if err != nil {
		return nil, err
	}
	body, err := NewParallelBodyFromMap(data, depth)
	if err != nil {
		return nil, err
	}
	state := &ParallelState{
		BaseState:    bs,
		ParallelBody: body,
		_Depth:       depth,
	}
	err = state.Init()
	if err != nil {
		return nil, err
	}
	return state, nil
}

// NewParallelBodyFromMap  NewParallelBodyFromMap
func NewParallelBodyFromMap(data map[string]interface{}, depth int) (*ParallelBody, error) {

	parallelbody := &ParallelBody{
		Branches:  []map[string]interface{}{},
		_Branches: []*StateMachineBody{},
	}
	err := InitParallelBodyByMap(parallelbody, data, depth)
	if err != nil {
		return nil, err
	}
	return parallelbody, nil
}

func InitParallelBodyByMap(body *ParallelBody, data map[string]interface{}, depth int) error {

	var err error

	err = DecodeMapToStruct(data, &body)
	if err != nil {
		return err
	}

	err = myvalidate.Struct(body)
	if err != nil {
		return err
	}
	if len(body.Branches) == 0 {
		return fmt.Errorf("%w:%s", vo.ErrorParallelBranchNumberLimitExceeded, "at least one branch, current branch number is 0")
	}
	if len(body.Branches) > ParserQuota.MaxParallelBranchNumber {
		return vo.ErrorParallelBranchNumberLimitExceeded
	}

	for _, branchdata := range body.Branches {
		branchsm, err := NewStateMachineBodyFromMap(branchdata, depth+1)
		if err != nil {
			return err
		}

		body._Branches = append(body._Branches, branchsm)
	}

	return nil
}

// GetBranches 返回多个分支的statemachine
func (s *ParallelState) GetBranches() []*StateMachineBody {
	return s._Branches
}

func (s *ParallelState) GetBaseState() *BaseState {
	return s.BaseState
}

func (s *ParallelState) SetDepth(depth int) {
	s._Depth = depth
}

func (s *ParallelState) GetDepth() int {
	return s._Depth
}

func (s *ParallelState) GetDefinition() (string, error) {
	data, err := ToString(s)
	return data, err
}
