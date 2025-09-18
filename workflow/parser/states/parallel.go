package states

import (
	"fmt"

	"gopkg.mihoyo.com/plat/cloudflow/workflow/quota"
	"gopkg.mihoyo.com/plat/cloudflow/workflow/vo"
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
	*BaseState
	Branches []map[string]interface{} `mapstructure:"Branches" validate:"required,gt=0"`
	// 这里修改成字典
	_Branches []*StateGroup
	_Depth    int
}

// NewParallelStateFromString NewParallelStateFromString
func NewParallelStateFromString(definition string, depth int) (*ParallelState, error) {

	state, err := NewGenericGroupStateFromString(definition, depth, NewParallelStateFromMap)
	return state, err
}

// NewParallelStateFromMap NewParallelStateFromMap
func NewParallelStateFromMap(data map[string]interface{}, depth int) (*ParallelState, error) {

	var err error
	state := &ParallelState{
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
func (s *ParallelState) InitByMap(data map[string]interface{}) error {
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

func (s *ParallelState) Init() error {

	s._Bone.Branches = []StateBone{}

	for idx, subbranch := range s.Branches {

		subsg, err := NewStateGroupFromMap(subbranch, s._Depth+1)
		if err != nil {
			return err
		}
		if subsg.GetName() == "" {
			subsg.SetName(fmt.Sprintf("%s__#%d", s.GetName(), idx+1))
		}
		err = s.AddBranchStateGroup(subsg)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddBranchStateGroup 添加一个分支
func (s *ParallelState) AddBranchStateGroup(sm *StateGroup) error {

	s._Branches = append(s._Branches, sm)
	smbone := sm.GetBone()
	s._Bone.Branches = append(s._Bone.Branches, smbone)

	// 限制并行分支数量
	if len(s._Branches) > quota.MaxParallelBranchNumber {
		return vo.ErrorParallelBranchNumberLimitExceeded
	}
	return nil
}

// GetBranches 返回多个分支的statemachine
func (s *ParallelState) GetBranches() []*StateGroup {
	return s._Branches
}

func (s *ParallelState) SetDepth(depth int) {
	s._Depth = depth
}
