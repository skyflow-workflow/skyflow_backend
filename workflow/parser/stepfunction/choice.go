package stepfunction

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/expression/stepfunction"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

type ChoiceBody struct {
	Choices []map[string]any `mapstructure:"Choices"`
	Default string           `mapstructure:"Default"`
}

// DecodeTaskState ...
func (decoder *StepfuncionDecoder) DecodeChoiceState(ctx context.Context, basestate *states.BaseState, data map[string]any) (
	states.State, error) {
	choicebody, err := decoder.DecodeChoiceBody(ctx, data)
	if err != nil {
		return nil, err
	}
	choice := &states.Choice{
		BaseState:  basestate,
		ChoiceBody: choicebody,
	}
	return choice, nil
}

// DecodeTaskBody ...
func (decoder *StepfuncionDecoder) DecodeChoiceBody(ctx context.Context, data map[string]any) (
	*states.ChoiceBody, error) {
	choicebody := ChoiceBody{}
	statebody := states.ChoiceBody{}
	err := decoder.MapDecode(data, &choicebody)
	if err != nil {
		return nil, err
	}
	statebody.Choices = make([]states.ChoiceBranch, len(choicebody.Choices))
	for idx, branchData := range choicebody.Choices {
		choiceBranch := states.ChoiceBranch{}
		err := decoder.MapDecode(branchData, &choiceBranch)
		if err != nil {
			return nil, err
		}
		condition, err := stepfunction.ParseEvaluateCondition(branchData)
		if err != nil {
			return nil, err
		}
		statebody.Choices[idx] = states.ChoiceBranch{
			Condition: condition,
			Next:      choiceBranch.Next,
		}
	}
	statebody.Default = choicebody.Default
	return &statebody, nil
}
