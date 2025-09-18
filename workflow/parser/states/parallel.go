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

// Parallel 并行状态
type Parallel struct {
	*BaseState
	*ParallelBody
}

func (p *Parallel) GetBaseState() *BaseState {
	return p.BaseState
}

type ParallelBody struct {
	Branches []map[string]interface{} `mapstructure:"Branches" validate:"required,gt=0"`
	// 这里修改成字典
	_Branches []*StateMachineBody
	_Depth    int
}

func (s *Parallel) Init() error {

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
func (s *Parallel) AddBranchStateGroup(smb *StateMachineBody) error {

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
func (s *Parallel) GetBranches() []*StateGroup {
	return s._Branches
}

func (s *Parallel) SetDepth(depth int) {
	s._Depth = depth
}
