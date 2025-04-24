// package expression show how to evaluate a boolean expression in workflow
package expression

import "github.com/skyflow-workflow/skyflow_backbend/workflow/expression/stepfunction"

// BooleanExpression is a interface for boolean expression evaluator
type BooleanExpression interface {
	// Evaluate evaluate the expression
	Evaluate(input any) bool
}

func NewStepfunctionExpression(expression map[string]any) (BooleanExpression, error) {
	stepfunction, err := stepfunction.NewStepFunctionExpression(expression)
	if err != nil {
		return nil, err
	}
	return stepfunction, nil
}
