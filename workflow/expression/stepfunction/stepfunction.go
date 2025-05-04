// stepfunction stepfunction style expression evaluator
package stepfunction

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/skyflow-workflow/skyflow_backbend/pkg/jsonpath"
)

// myValidate self define validator
var myValidate = validator.New()

func init() {
	value := reflect.ValueOf(Operators)
	for i := 0; i < value.NumField(); i++ {
		ChoiceOperatorArray = append(ChoiceOperatorArray, value.Field(i).String())
	}
}

var KeyVariable = "Variable"

// Operators compare operators
var Operators = struct {
	// exist  is variable exist
	VariableExist string
	IsExist       string
	IsPresent     string
	// Bool  is boolean
	IsBoolean   string
	IsNull      string
	IsNumeric   string
	IsString    string
	IsTimestamp string
	// String compare
	StringEquals            string
	StringLessThan          string
	StringGreaterThan       string
	StringLessThanEquals    string
	StringGreaterThanEquals string
	// Numeric compare
	NumericEquals            string
	NumericLessThan          string
	NumericGreaterThan       string
	NumericLessThanEquals    string
	NumericGreaterThanEquals string
	NumericEqualsPath        string
	// bool compare
	BooleanEquals string
	// Path string
	BooleanEqualsPath string
	// Timestamp compare
	TimestampEquals            string
	TimestampLessThan          string
	TimestampGreaterThan       string
	TimestampLessThanEquals    string
	TimestampGreaterThanEquals string
}{
	// exist is variable exist
	VariableExist: "VariableExist",
	IsExist:       "IsExist",
	IsPresent:     "IsPresent",
	// Bool is variable is specified type
	IsBoolean:   "IsBoolean",
	IsNull:      "IsNull",
	IsNumeric:   "IsNumeric",
	IsString:    "IsString",
	IsTimestamp: "IsTimestamp",
	// String string compare
	StringEquals:            "StringEquals",
	StringLessThan:          "StringLessThan",
	StringGreaterThan:       "StringGreaterThan",
	StringLessThanEquals:    "StringLessThanEquals",
	StringGreaterThanEquals: "StringGreaterThanEquals",
	// Numeric numeric compare
	NumericEquals:            "NumericEquals",
	NumericLessThan:          "NumericLessThan",
	NumericGreaterThan:       "NumericGreaterThan",
	NumericLessThanEquals:    "NumericLessThanEquals",
	NumericGreaterThanEquals: "NumericGreaterThanEquals",
	NumericEqualsPath:        "NumericEqualsPath",
	// Boolean boolean compare
	BooleanEquals: "BooleanEquals",
	// Path path compare
	BooleanEqualsPath: "BooleanEqualsPath",
	// Timestamp timestamp compare
	TimestampEquals:            "TimestampEquals",
	TimestampLessThan:          "TimestampLessThan",
	TimestampGreaterThan:       "TimestampGreaterThan",
	TimestampLessThanEquals:    "TimestampLessThanEquals",
	TimestampGreaterThanEquals: "TimestampGreaterThanEquals",
}

// CompareCode compare operators
var CompareCode = struct {
	GT  string
	GTE string
	EQ  string
	LT  string
	LTE string
}{
	GT:  "GT",
	GTE: "GTE",
	EQ:  "EQ",
	LT:  "LT",
	LTE: "LTE",
}

// CompareFuncs compare functions map
// check if the current value is match to the except value
// @param current is input variable
// @param except is evaluate expect value
var CompareFuncs = map[string]func(current any, except any) bool{
	Operators.IsBoolean: func(current any, except any) bool {
		return IsBool(current) == except.(bool)
	},
	Operators.IsNull: func(current any, except any) bool {
		return IsNull(current) == except.(bool)
	},
	Operators.IsNumeric: func(current any, except any) bool {
		return IsNumeric(current) == except.(bool)
	},
	Operators.IsString: func(current any, except any) bool {
		return IsString(current) == except.(bool)
	},
	// Timestamp
	Operators.IsTimestamp: func(current any, except any) bool {
		return IsTimestamp(current) == except.(bool)
	},
	Operators.StringEquals: func(current any, except any) bool {
		return StringCompare(current, except, CompareCode.EQ)
	},
	Operators.StringLessThan: func(current any, except any) bool {
		return StringCompare(current, except, CompareCode.LT)
	},
	Operators.StringLessThanEquals: func(current any, except any) bool {
		return StringCompare(current, except, CompareCode.LTE)
	},
	Operators.StringGreaterThan: func(current any, except any) bool {
		return StringCompare(current, except, CompareCode.GT)
	},
	Operators.StringGreaterThanEquals: func(current any, except any) bool {
		return StringCompare(current, except, CompareCode.GTE)
	},
	// NumericEquals
	Operators.NumericEquals: func(current any, except any) bool {
		return NumberCompare(current, except, CompareCode.EQ)
	},
	Operators.NumericLessThan: func(current any, except any) bool {
		return NumberCompare(current, except, CompareCode.LT)
	},
	Operators.NumericGreaterThan: func(current any, except any) bool {
		return NumberCompare(current, except, CompareCode.GT)
	},
	Operators.NumericLessThanEquals: func(current any, except any) bool {
		return NumberCompare(current, except, CompareCode.LTE)
	},
	Operators.NumericGreaterThanEquals: func(current any, except any) bool {
		return NumberCompare(current, except, CompareCode.GTE)
	},

	// bool
	Operators.BooleanEquals: func(current any, except any) bool {
		return current.(bool) == except.(bool)
	},
}

// ExistFuncs check if the variable is exist
// return true if the variable match the expect result
// @param current is input variable is exist
// @param except is evaluate expect result
var ExistFuncs = map[string]func(current bool, except bool) bool{
	Operators.IsPresent:     func(current bool, except bool) bool { return current == except },
	Operators.IsExist:       func(current bool, except bool) bool { return current == except },
	Operators.VariableExist: func(current bool, except bool) bool { return current == except },
}

// ChoiceOperatorArray choice operator list
var ChoiceOperatorArray []string

// IsString check if the value is string
func IsString(i any) bool {
	return reflect.TypeOf(i).Kind() == reflect.String
}

// IsBool check if the value is bool
func IsBool(i any) bool {
	return reflect.TypeOf(i).Kind() == reflect.Bool
}

// IsNull check if the value is nil
func IsNull(i any) bool {
	return i == nil
}

// IsTimestamp check if the value is timestamp
func IsTimestamp(i any) bool {
	return false
}

// IsNumeric check if the value is numeric
func IsNumeric(i any) bool {
	t := reflect.TypeOf(i).Kind()
	if t == reflect.Float64 || t == reflect.Float32 || t == reflect.Int ||
		t == reflect.Int64 {
		return true
	}
	return false
}

// ToNumber convert to number float64
func ToNumber(i any) (float64, bool) {
	var v float64
	t := reflect.TypeOf(i).Kind()
	switch t {
	case reflect.Float64:
		v = i.(float64)
	case reflect.Float32:
		v = float64(i.(float32))
	case reflect.Int:
		v = float64(i.(int))
	case reflect.Int32:
		v = float64(i.(int32))
	case reflect.Int64:
		v = float64(i.(int64))

	default:
		return v, false
	}
	return v, true
}

// NumberCompare compare number
// @param check is input variable
// @param except is evaluate expect value
// @param op is compare operator
func NumberCompare(check any, except any, op string) bool {

	var checkvalue float64
	var exceptvalue float64
	var ok bool
	if checkvalue, ok = ToNumber(check); !ok {
		return false
	}

	if exceptvalue, ok = ToNumber(except); !ok {
		return false
	}

	switch op {
	case CompareCode.GT:
		return checkvalue > exceptvalue
	case CompareCode.GTE:
		return checkvalue >= exceptvalue
	case CompareCode.EQ:
		return checkvalue == exceptvalue
	case CompareCode.LT:
		return checkvalue < exceptvalue
	case CompareCode.LTE:
		return checkvalue <= exceptvalue
	}
	return false
}

// StringCompare compare string
// @param check is input variable
// @param except is evaluate expect value
// @param op is compare operator
func StringCompare(check any, except any, op string) bool {

	var checkvalue string
	var exceptvalue string
	var ok bool
	if checkvalue, ok = check.(string); !ok {
		return false
	}
	if exceptvalue, ok = except.(string); !ok {
		return false
	}
	switch op {
	case CompareCode.GT:
		return checkvalue > exceptvalue
	case CompareCode.GTE:
		return checkvalue >= exceptvalue
	case CompareCode.EQ:
		return checkvalue == exceptvalue
	case CompareCode.LT:
		return checkvalue < exceptvalue
	case CompareCode.LTE:
		return checkvalue <= exceptvalue
	}
	return false
}

// TimeStampCompare compare timestamp
func TimeStampCompare(check any, except any, op string) bool {

	return true
}

// ChoiceSetOperator choice set operator
var ChoiceSetOperator = struct {
	And  string
	Not  string
	Or   string
	None string
}{
	And:  "And",
	Not:  "Not",
	Or:   "Or",
	None: "None",
}

// EvaluateUnit evaluate unit
type EvaluateUnit struct {
	// compare operate e.g. 'NumberEqual'
	Operator string
	// compare value e.g  '12'
	Operand any
	// variable path e.g  '$.foo'
	Variable string
}

func (eu *EvaluateUnit) Evaluate(input any) bool {
	var err error
	var checkvalue any
	checkvalue, err = jsonpath.JsonPathGetValue(eu.Variable, input)

	// check if the variable is exist first
	// if the err is nil, the variable is exist
	// if the err is not nil, the variable is not exist
	if evaluateFunc, ok := ExistFuncs[eu.Operator]; ok {
		exist := evaluateFunc(err == nil, eu.Operand.(bool))
		return exist
	}
	// if the variable is not exist, the evaluate result is false
	if err != nil {
		return false
	}

	evaluateFunc, ok := CompareFuncs[eu.Operator]
	if !ok {
		return false
	}
	match := evaluateFunc(checkvalue, eu.Operand)
	return match
}

// EvaluateCondition evaluate condition, combine evaluate unit with And/Or/Not
type EvaluateCondition struct {
	And []EvaluateCondition
	Not *EvaluateCondition
	Or  []EvaluateCondition
	*EvaluateUnit
}

func (ec *EvaluateCondition) Evaluate(input any) bool {
	return ec.Validate(input)
}

func (ec *EvaluateCondition) Validate(input any) bool {

	if ec.Not != nil {
		return !ec.Not.Validate(input)
	} else if len(ec.Or) > 0 {
		for _, subEc := range ec.Or {
			if subEc.Validate(input) {
				return true
			}
		}
		return false
	} else if len(ec.And) > 0 {
		for _, subEc := range ec.And {
			if !subEc.Validate(input) {
				return false
			}
		}
		return true
	} else if ec.EvaluateUnit != nil {
		return ec.EvaluateUnit.Evaluate(input)
	}
	return false
}

// StepFunctionExpression stepfunction style expression evaluator
type StepFunctionExpression struct {
	Expression *EvaluateCondition
}

// NewStepFunctionExpression create a new stepfunction expression
func NewStepFunctionExpression(expression map[string]any) (*StepFunctionExpression, error) {
	var err error
	cond, err := ParseEvaluateCondition(expression)
	if err != nil {
		return nil, err
	}
	return &StepFunctionExpression{
		Expression: cond,
	}, nil
}

// Evaluate evaluate the expression
func (e *StepFunctionExpression) Evaluate(input any) bool {
	return e.Expression.Evaluate(input)
}

type InputEvaluateSetUnit struct {
	And []map[string]any `mapstructure:"And"`
	Not map[string]any   `mapstructure:"Not"`
	Or  []map[string]any `mapstructure:"Or"`
}

// InputEvaluateUnit 基础比较单元
type InputEvaluateUnit struct {
	Variable string `mapstructure:"Variable" validate:"required,gt=0,startswith=$"`
	// Exist
	VariableExist bool `mapstructure:"VariableExist"`
	IsExist       bool
	IsPresent     bool
	// bool 比较
	IsBoolean   bool `mapstructure:"IsBoolean"`
	IsNull      bool `mapstructure:"IsNull"`
	IsNumeric   bool `mapstructure:"IsNumeric"`
	IsString    bool `mapstructure:"IsString"`
	IsTimestamp bool `mapstructure:"IsTimestamp"`
	// value 比较
	StringEquals               string  `mapstructure:"StringEquals"`
	StringLessThan             string  `mapstructure:"StringLessThan"`
	StringGreaterThan          string  `mapstructure:"StringGreaterThan"`
	StringLessThanEquals       string  `mapstructure:"StringLessThanEquals"`
	StringGreaterThanEquals    string  `mapstructure:"StringGreaterThanEquals"`
	NumericEquals              float64 `mapstructure:"NumericEquals"`
	NumericLessThan            float64 `mapstructure:"NumericLessThan"`
	NumericGreaterThan         float64 `mapstructure:"NumericGreaterThan"`
	NumericLessThanEquals      float64 `mapstructure:"NumericLessThanEquals"`
	NumericGreaterThanEquals   float64 `mapstructure:"NumericGreaterThanEquals"`
	TimestampEquals            string  `mapstructure:"TimestampEquals"`
	TimestampLessThan          string  `mapstructure:"TimestampLessThan"`
	TimestampGreaterThan       string  `mapstructure:"TimestampGreaterThan"`
	TimestampLessThanEquals    string  `mapstructure:"TimestampLessThanEquals"`
	TimestampGreaterThanEquals string  `mapstructure:"TimestampGreaterThanEquals"`
	BooleanEquals              bool    `mapstructure:"BooleanEquals"`
}

func ParseEvaluateCondition(data map[string]any) (*EvaluateCondition, error) {

	var err error
	resp := &EvaluateCondition{}
	evaluateSetUnit := &InputEvaluateSetUnit{}

	err = mapstructure.Decode(data, evaluateSetUnit)
	if err != nil {
		return nil, err
	}

	if len(evaluateSetUnit.And) > 0 {
		for _, subData := range evaluateSetUnit.And {
			subSetUnit, err := ParseEvaluateCondition(subData)
			if err != nil {
				return nil, err
			}
			resp.And = append(resp.And, *subSetUnit)
		}
	} else if len(evaluateSetUnit.Or) > 0 {
		for _, subData := range evaluateSetUnit.Or {
			subSetUnit, err := ParseEvaluateCondition(subData)
			if err != nil {
				return nil, err
			}
			resp.Or = append(resp.Or, *subSetUnit)
		}
	} else if evaluateSetUnit.Not != nil {
		subSetUnit, err := ParseEvaluateCondition(evaluateSetUnit.Not)
		if err != nil {
			return nil, err
		}
		resp.Not = subSetUnit
	} else {
		evaluateUnit, err := ParseEvaluateUnit(data)
		if err != nil {
			return nil, err
		}
		resp.EvaluateUnit = evaluateUnit
	}
	return resp, nil
}

func ParseEvaluateUnit(data map[string]any) (*EvaluateUnit, error) {

	var inputEvaluateUnit InputEvaluateUnit
	var evaluateUnit EvaluateUnit
	var err error
	// 初始化自身
	err = mapstructure.Decode(data, &inputEvaluateUnit)
	if err != nil {
		return nil, err
	}
	err = myValidate.Struct(inputEvaluateUnit)
	if err != nil {
		return nil, err
	}

	//  验证是否是 jsonpath
	_, err = jsonpath.JsonPathCompile(inputEvaluateUnit.Variable)
	if err != nil {
		err = fmt.Errorf("key '%s' jsonpath compile error: %w", KeyVariable, err)
		return nil, err
	}
	evaluateUnit.Variable = inputEvaluateUnit.Variable

	var opFind = false
	for _, operator := range ChoiceOperatorArray {
		if value, ok := data[operator]; ok {
			evaluateUnit.Operator = operator
			evaluateUnit.Operand = value
			opFind = true
			break
		}
	}
	if !opFind {
		err = fmt.Errorf("missing operate expression")
		return nil, err
	}
	return &evaluateUnit, nil
}
