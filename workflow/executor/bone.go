package executor

import (
	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
)

// Bone is the base struct of ExecutionStep, Bone is used for draw the workflow graph.

type BaseBone struct {
	states.BaseBone
	// Status of the step
	Status string `json:"Status"`
	// StepID is the id of the step
	StepID int `json:"StepID"`
	// Index is the index of the step in the execution
	Index int `json:"Index"`
}

type StepBone struct {
	BaseBone
	*ExecutionBone
	// for map /parallel contains multiple branches
	Branches []StepBone `json:"Branches"`
}

// ExecutionBone  Execution Bone
type ExecutionBone struct {
	Version string              `json:"Version"`
	StartAt string              `json:"StartAt"`
	States  map[string]StepBone `json:"States"`
}
