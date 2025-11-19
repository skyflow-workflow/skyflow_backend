package executor

import "github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"

// Bone is the base struct of ExecutionStep, Bone is used for draw the workflow graph.

type BaseBone struct {
	states.BaseBone
	Status string `json:"Status"`
	StepID int    `json:"StepID"`
	Index  int    `json:"Index"`
}

type StepBone struct {
	BaseBone
	*StateMachineBone
	// for map /parallel contains multiple branches
	Branches []StepBone `json:"Branches"`
}

// ExecutionBone  Execution Bone
type ExecutionBone struct {
	StartAt string              `json:"StartAt"`
	States  map[string]StepBone `json:"States"`
}

type StateMachineBone struct {
	// Type    string
	Version string
	StartAt string              `json:"StartAt"`
	States  map[string]StepBone `json:"States"`
}
