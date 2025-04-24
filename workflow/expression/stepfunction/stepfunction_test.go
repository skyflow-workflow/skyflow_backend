package stepfunction

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mitchellh/mapstructure"
)

func TestEvaluateUnit(t *testing.T) {
	var testcases = []struct {
		name   string
		unit   EvaluateUnit
		input  map[string]any
		expect bool
	}{
		{
			name: "single variable existed",
			unit: EvaluateUnit{
				Operator: "VariableExist",
				Operand:  true,
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": "value1", "key2": "value2"},
			expect: true,
		},
		{
			name: "single variable not existed",
			unit: EvaluateUnit{
				Operator: "VariableExist",
				Operand:  false,
				Variable: "$.key3",
			},
			input:  map[string]any{"key1": "value1", "key2": "value2"},
			expect: true,
		},
		{
			name: "single variable not existed",
			unit: EvaluateUnit{
				Operator: "VariableExist",
				Operand:  true,
				Variable: "$.key3",
			},
			input:  map[string]any{"key1": "value1", "key2": "value2"},
			expect: false,
		},
		{
			name: "single variable not existed for nil input",
			unit: EvaluateUnit{
				Operator: "VariableExist",
				Operand:  false,
				Variable: "$.key3",
			},
			input:  nil,
			expect: true,
		},
		{
			name: "single variable for string comparison equal type error",
			unit: EvaluateUnit{
				Operator: "StringEquals",
				Operand:  111,
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": "value1", "key2": "value2"},
			expect: false,
		},
		{
			name: "single variable for string comparison equal",
			unit: EvaluateUnit{
				Operator: "StringEquals",
				Operand:  "value1",
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": "value1", "key2": "value2"},
			expect: true,
		},
		{
			name: "single variable for string comparison not equal",
			unit: EvaluateUnit{
				Operator: "StringEquals",
				Operand:  "value1",
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": "value2", "key2": "value2"},
			expect: false,
		},
		{
			name: "variable for string comparison greater than",
			unit: EvaluateUnit{
				Operator: "StringGreaterThan",
				Operand:  "value1",
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": "value2", "key2": "value2"},
			expect: true,
		},
		{
			name: "variable for string comparison less than or equal",
			unit: EvaluateUnit{
				Operator: "StringLessThanEquals",
				Operand:  "value2",
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": "value1", "key2": "value2"},
			expect: true,
		},
		{
			name: "variable for string comparison not less than",
			unit: EvaluateUnit{
				Operator: "StringLessThan",
				Operand:  "value1",
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": "value1", "key2": "value2"},
			expect: false,
		},
		{
			name: "variable for string comparison not greater than",
			unit: EvaluateUnit{
				Operator: "StringGreaterThan",
				Operand:  "value1",
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": "value1", "key2": "value2"},
			expect: false,
		},
		{
			name: "variable for string comparison not greater than",
			unit: EvaluateUnit{
				Operator: "StringGreaterThan",
				Operand:  "value1",
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": "value1", "key2": "value2"},
			expect: false,
		},
		{
			name: "variable for numeric comparison not numeric",
			unit: EvaluateUnit{
				Operator: "NumericGreaterThan",
				Operand:  10,
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": "value1", "key2": "value2"},
			expect: false,
		},
		{
			name: "variable for numeric comparison equals",
			unit: EvaluateUnit{
				Operator: "NumericEquals",
				Operand:  10,
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": 10, "key2": 10},
			expect: true,
		},
		{
			name: "variable for numeric comparison equals float",
			unit: EvaluateUnit{
				Operator: "NumericEquals",
				Operand:  10.1,
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": 10.1, "key2": 10.0},
			expect: true,
		},
		{
			name: "variable for numeric comparison greater than",
			unit: EvaluateUnit{
				Operator: "NumericGreaterThan",
				Operand:  10,
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": 11, "key2": 10},
			expect: true,
		},
		{
			name: "variable for numeric comparison less than",
			unit: EvaluateUnit{
				Operator: "NumericLessThan",
				Operand:  10,
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": 9, "key2": 10},
			expect: true,
		},
		{
			name: "variable for numeric comparison less than or equal",
			unit: EvaluateUnit{
				Operator: "NumericLessThanEquals",
				Operand:  10,
				Variable: "$.key1",
			},
			input:  map[string]any{"key1": 10, "key2": 10},
			expect: true,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			exist := testcase.unit.Evaluate(testcase.input)
			assert.Equal(t, exist, testcase.expect)
		})
	}
}

func TestParseEvaluateCondition(t *testing.T) {
	var testcases = []struct {
		name        string
		expression  map[string]any
		expect      *EvaluateCondition
		expectError error
	}{
		{
			name: "single variable existed",
			expression: map[string]any{
				"Variable":      "$.key1",
				"VariableExist": true,
			},
			expect: &EvaluateCondition{
				EvaluateUnit: &EvaluateUnit{
					Operator: "VariableExist",
					Operand:  true,
					Variable: "$.key1",
				},
			},
		},
		{
			name: "and condition",
			expression: map[string]any{
				"And": []map[string]any{
					{
						"StringEquals": "value1",
						"Variable":     "$.key1",
					},
					{
						"StringEquals": "value2",
						"Variable":     "$.key2",
					},
				},
			},
			expect: &EvaluateCondition{
				And: []EvaluateCondition{
					{
						EvaluateUnit: &EvaluateUnit{
							Operator: "StringEquals",
							Operand:  "value1",
							Variable: "$.key1",
						},
					},
					{
						EvaluateUnit: &EvaluateUnit{
							Operator: "StringEquals",
							Operand:  "value2",
							Variable: "$.key2",
						},
					},
				},
			},
		},
		{
			name: "or condition",
			expression: map[string]any{
				"Or": []map[string]any{
					{
						"StringEquals": "value1",
						"Variable":     "$.key1",
					},
					{
						"NumericEquals": 10,
						"Variable":      "$.key2",
					},
				},
			},
			expect: &EvaluateCondition{
				Or: []EvaluateCondition{
					{
						EvaluateUnit: &EvaluateUnit{
							Operator: "StringEquals",
							Operand:  "value1",
							Variable: "$.key1",
						},
					},
					{
						EvaluateUnit: &EvaluateUnit{
							Operator: "NumericEquals",
							Operand:  10,
							Variable: "$.key2",
						},
					},
				},
			},
		},
		{
			name: "not condition",
			expression: map[string]any{
				"Not": map[string]any{
					"StringEquals": "value1",
					"Variable":     "$.key1",
				},
			},
			expect: &EvaluateCondition{
				Not: &EvaluateCondition{
					EvaluateUnit: &EvaluateUnit{
						Operator: "StringEquals",
						Operand:  "value1",
						Variable: "$.key1",
					},
				},
			},
		},
		{
			name: "two level condition",
			expression: map[string]any{
				"And": []map[string]any{
					{
						"Or": []map[string]any{
							{
								"StringEquals": "value1",
								"Variable":     "$.key1",
							},
							{
								"StringEquals": "value2",
								"Variable":     "$.key2",
							},
						},
					},
					{
						"Or": []map[string]any{
							{
								"StringEquals": "value3",
								"Variable":     "$.key3",
							},
							{
								"StringEquals": "value4",
								"Variable":     "$.key4",
							},
						},
					},
				},
			},
			expect: &EvaluateCondition{
				And: []EvaluateCondition{
					{
						Or: []EvaluateCondition{
							{
								EvaluateUnit: &EvaluateUnit{
									Operator: "StringEquals",
									Operand:  "value1",
									Variable: "$.key1",
								},
							},
							{
								EvaluateUnit: &EvaluateUnit{
									Operator: "StringEquals",
									Operand:  "value2",
									Variable: "$.key2",
								},
							},
						},
					},
					{
						Or: []EvaluateCondition{
							{
								EvaluateUnit: &EvaluateUnit{
									Operator: "StringEquals",
									Operand:  "value3",
									Variable: "$.key3",
								},
							},
							{
								EvaluateUnit: &EvaluateUnit{
									Operator: "StringEquals",
									Operand:  "value4",
									Variable: "$.key4",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			parsed, err := ParseEvaluateCondition(testcase.expression)
			if err != nil {
				assert.Equal(t, testcase.expectError, nil)
				assert.Equal(t, testcase.expectError.Error(), err.Error())
				return
			}
			t.Log(parsed)
			assert.Equal(t, parsed, testcase.expect)
		})
	}
}

func TestParseEvaluateUnit(t *testing.T) {
	var testcases = []struct {
		name        string
		expression  map[string]any
		expect      *EvaluateUnit
		expectError error
	}{
		{
			name: "variable existed",
			expression: map[string]any{
				"Variable":      "$.key1",
				"VariableExist": true,
			},
			expect: &EvaluateUnit{
				Operator: "VariableExist",
				Operand:  true,
				Variable: "$.key1",
			},
		},
		{
			name: "string equals operator",
			expression: map[string]any{
				"StringEquals": "value1",
				"Variable":     "$.key1",
			},
			expect: &EvaluateUnit{
				Operator: "StringEquals",
				Operand:  "value1",
				Variable: "$.key1",
			},
		},
		{
			name: "numeric equals operator",
			expression: map[string]any{
				"NumericEquals": 10,
				"Variable":      "$.key1",
			},
			expect: &EvaluateUnit{
				Operator: "NumericEquals",
				Operand:  10,
				Variable: "$.key1",
			},
		},
		{
			name: "boolean equals operator",
			expression: map[string]any{
				"BooleanEquals": false,
				"Variable":      "$.key1",
			},
			expect: &EvaluateUnit{
				Operator: "BooleanEquals",
				Operand:  false,
				Variable: "$.key1",
			},
		},
		{
			name: "multiple operators",
			expression: map[string]any{
				"BooleanEquals": false,
				"Variable":      "$.key1",
				"StringEquals":  "value1",
			},
			expect: &EvaluateUnit{
				Operator: "StringEquals",
				Operand:  "value1",
				Variable: "$.key1",
			},
		},
		{
			name: "missing operator",
			expression: map[string]any{
				"Variable": "$.key1",
			},
			expectError: fmt.Errorf("missing operate expression"),
		},
		{
			name: "wrong operator value",
			expression: map[string]any{
				"Variable":     "$.key1",
				"StringEquals": false,
			},
			expectError: &mapstructure.Error{
				Errors: []string{
					"'StringEquals' expected type 'string', got unconvertible type 'bool', value: 'false'",
				},
			},
		},
		{
			name: "compile jsonpath error",
			expression: map[string]any{
				"Variable":     "$key1",
				"StringEquals": "value1",
			},
			expectError: fmt.Errorf("key 'Variable' jsonpath compile error: parse error at 2 in $key1"),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			parsed, err := ParseEvaluateUnit(testcase.expression)
			if err != nil {
				targetError := &mapstructure.Error{}
				if errors.As(err, &targetError) {
					assert.Equal(t, testcase.expectError.Error(), err.Error())
					return
				}
				assert.Equal(t, err.Error(), testcase.expectError.Error())
				return
			}
			assert.Equal(t, testcase.expectError, nil)
			t.Log(parsed)
			assert.Equal(t, parsed, testcase.expect)
		})
	}
}

func TestStepFunctionExpression(t *testing.T) {
	var testcases = []struct {
		name        string
		expression  map[string]any
		input       map[string]any
		expect      bool
		expectError error
	}{
		{
			name: "single variable existed",
			expression: map[string]any{
				"VariableExist": true,
				"Variable":      "$.foo",
			},
			input: map[string]any{
				"foo": "bar",
			},
			expect:      true,
			expectError: nil,
		},
		{
			name: "single variable not existed",
			expression: map[string]any{
				"VariableExist": false,
				"Variable":      "$.foo",
			},
			input: map[string]any{
				"foo": "bar",
			},
			expect:      false,
			expectError: nil,
		},
		{
			name: "single variable not existed for nil input",
			expression: map[string]any{
				"VariableExist": false,
				"Variable":      "$.foo",
			},
			input:       nil,
			expect:      true,
			expectError: nil,
		},
		{
			name: "string equals operator",
			expression: map[string]any{
				"StringEquals": "value1",
				"Variable":     "$.key1",
			},
			input: map[string]any{
				"key1": "value1",
			},
			expect:      true,
			expectError: nil,
		},
		{
			name: "string not equals operator",
			expression: map[string]any{
				"StringEquals": "value1",
				"Variable":     "$.key1",
			},
			input: map[string]any{
				"key1": "value2",
			},
			expect:      false,
			expectError: nil,
		},
		{
			name: "numeric equals operator",
			expression: map[string]any{
				"NumericEquals": 10,
				"Variable":      "$.key1",
			},
			input: map[string]any{
				"key1": 10,
			},
			expect:      true,
			expectError: nil,
		},
		{
			name: "numeric not equals operator",
			expression: map[string]any{
				"NumericEquals": 10,
				"Variable":      "$.key1",
			},
			input: map[string]any{
				"key1": 11,
			},
			expect:      false,
			expectError: nil,
		},
		{
			name: "boolean equals operator",
			expression: map[string]any{
				"BooleanEquals": false,
				"Variable":      "$.key1",
			},
			input: map[string]any{
				"key1": false,
			},
			expect:      true,
			expectError: nil,
		},
		{
			name: "boolean not equals operator",
			expression: map[string]any{
				"BooleanEquals": false,
				"Variable":      "$.key1",
			},
			input: map[string]any{
				"key1": true,
			},
			expect:      false,
			expectError: nil,
		},
		{
			name: "two level condition with and",
			expression: map[string]any{
				"And": []map[string]any{
					{
						"BooleanEquals": false,
						"Variable":      "$.key1",
					},
					{
						"StringEquals": "value1",
						"Variable":     "$.key2",
					},
				},
			},
			input: map[string]any{
				"key1": false,
				"key2": "value1",
			},
			expect:      true,
			expectError: nil,
		},
		{
			name: "two level condition with or",
			expression: map[string]any{
				"Or": []map[string]any{
					{
						"BooleanEquals": true,
						"Variable":      "$.key1",
					},
					{
						"StringEquals": "value1",
						"Variable":     "$.key2",
					},
				},
			},
			input: map[string]any{
				"key1": false,
				"key2": "value1",
			},
			expect:      true,
			expectError: nil,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			exp, err := NewStepFunctionExpression(testcase.expression)
			if err != nil {
				assert.Equal(t, testcase.expectError, nil)
				assert.Equal(t, testcase.expectError.Error(), err.Error())
				return
			}
			assert.Equal(t, exp.Evaluate(testcase.input), testcase.expect)
		})
	}
}
