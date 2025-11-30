package states

import (
	"encoding/json"
	"fmt"

	"github.com/skyflow-workflow/skyflow_backend/workflow/expression"
)

// ChoiceFieldNames
var ChoiceFieldNames = struct {
	Choices string
	Default string
}{
	Choices: "Choices",
	Default: "Default",
}

type ChoiceBranch struct {
	Condition expression.BooleanExpression
	Next      string `mapstructure:"Next" validate:"required,gt=0"`
}

// ChoiceBody ...
type ChoiceBody struct {
	Choices []ChoiceBranch `mapstructure:"Choices"`
	Default string         `mapstructure:"Default" validate:"gte=0"`
}

// ChoiceState ...
type ChoiceState struct {
	*BaseState  `json:",inline"`
	*ChoiceBody `json:",inline"`
}

// NewChoiceStateFromString NewChoiceStateFromString
func NewChoiceStateFromString(definition string) (state *ChoiceState, err error) {

	var data = map[string]interface{}{}
	err = json.Unmarshal([]byte(definition), &data)
	if err != nil {
		return
	}
	state, err = NewChoiceStateFromMap(data)
	return
}

// NewChoiceStateFromMap NewChoiceStateFromMap
func NewChoiceStateFromMap(data map[string]interface{}) (state *ChoiceState, err error) {

	// basestate
	bs, err := NewBaseStateFromMap(data)
	if err != nil {
		return state, err
	}
	choicebody, err := NewChoiceBodyFromMap(data)
	if err != nil {
		return state, err
	}
	state = &ChoiceState{
		BaseState:  bs,
		ChoiceBody: choicebody,
	}
	return
}

func NewChoiceBodyFromMap(data map[string]interface{}) (*ChoiceBody, error) {
	var err error
	choicebody := &ChoiceBody{}
	err = InitChoiceBodyByMap(choicebody, data)
	if err != nil {
		return nil, err
	}
	return choicebody, nil

}

// InitByMap Initialize ChoiceState Content
func InitChoiceBodyByMap(body *ChoiceBody, data map[string]interface{}) error {

	var err error
	// 初始化自身
	err = DecodeMapToStruct(data, body)
	if err != nil {
		return err
	}
	err = myValidate.Struct(body)
	if err != nil {
		return err
	}

	if len(body.Choices) == 0 {
		err = fmt.Errorf("choice branch is empty")
		return err
	}

	choicebranchs, ok := data[ChoiceFieldNames.Choices].([]interface{})
	if !ok {
		err = fmt.Errorf("choice branch is not array")
		return err
	}
	for index, branch := range choicebranchs {
		branchmap, ok := branch.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("choice branch is not map")
			return err
		}
		exp, err := expression.NewStepfunctionExpression(branchmap)

		if err != nil {
			err = fmt.Errorf("choice branch index '%d': %w", index, err)
			return err
		}

		body.Choices[index].Condition = exp
	}
	return nil
}

func (choice *ChoiceState) GetBaseState() *BaseState {
	return choice.BaseState
}

// GetBone get choice bone
func (choice *ChoiceState) GetBone() StateBone {
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
func (choice *ChoiceState) GetNextState(input any) (NextState, error) {
	var err error
	var ns NextState
	newinput, err := choice.GetInput(input)
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
func (choice *ChoiceState) ChoiceNextState(input any) string {
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
func (choice *ChoiceState) IsEnd() bool {
	return false
}

func (choice *ChoiceState) GetDefinition() (string, error) {
	data, err := ToString(choice)
	return data, err
}
