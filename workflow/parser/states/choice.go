package states

import (
	"fmt"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/expression"
)

type ChoiceBranch struct {
	Condition expression.BooleanExpression
	Next      string `mapstructure:"Next" validate:"required,gt=0"`
}

// ChoiceBody ...
type ChoiceBody struct {
	Choices []ChoiceBranch `mapstructure:"Choices"`
	Default string         `mapstructure:"Default" validate:"gte=0"`
}

// Choice ...
type Choice struct {
	*BaseState
	*ChoiceBody
}

// GetBone get choice bone
func (choice *Choice) GetBone() StateBone {
	bone := choice.BaseState.GetBone()
	for _, choicebranch := range choice.Choices {
		bone.Next = append(bone.Next, choicebranch.Next)
	}
	if choice.Default != "" {
		bone.Next = append(bone.Next, choice.Default)
	}
	bone.End = false
	return bone
}

// GetNextState Get Next State
// input state origin input
func (choice *Choice) GetNextState(input any) (NextState, error) {
	var err error
	var ns NextState
	newinput, err := choice.GenParameters(input)
	if err != nil {
		return ns, err
	}
	next := choice.ChoiceNextState(newinput)
	if next == "" {
		err = fmt.Errorf("no choice branch match")
		return ns, err
	}

	// choice 应该是直接使用 input的
	ns = NextState{
		Name:   next,
		Output: input,
	}
	return ns, nil
}

// ChoiceNextState
func (choice *Choice) ChoiceNextState(input any) string {
	var success bool
	for _, branch := range choice.Choices {
		success = branch.Condition.Evaluate(input)
		if success {
			return branch.Next
		}
	}
	if choice.Default != "" {
		return choice.Default
	}
	return ""
}

// IsEnd check if choice is end, return false always, choice is not end state
func (choice *Choice) IsEnd() bool {
	return false
}
