package jsonpath

import (
	"github.com/ohler55/ojg/jp"
	"github.com/oliveagle/jsonpath"
)

// JsonPathSetValue sets a value in a JSONPath expression
//
// path is the JSONPath expression
// source is the JSON object to set the value in
// value is the value to set
func JsonPathSetValue(path string, source any, value any) error {

	jpExpr, err := jp.ParseString(path)
	if err != nil {
		return err
	}
	err = jpExpr.Set(source, value)
	if err != nil {
		return err
	}
	return nil
}

// JsonPathGetValue gets a value from a JSONPath expression
//
// path is the JSONPath expression
// source is the JSON object to get the value from
// returns the value at the path
func JsonPathGetValue(path string, source any) (any, error) {

	jpc, err := jsonpath.Compile(path)
	if err != nil {
		return nil, err
	}
	result, err := jpc.Lookup(source)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// JsonPathCompile compiles a JSONPath expression
//
// jsonpath is the JSONPath expression
// returns the compiled JSONPath expression
func JsonPathCompile(jsonpath string) (jp.Expr, error) {

	return jp.ParseString(jsonpath)
}
