package states

import "github.com/ohler55/ojg/jp"

// JSONPathCompiled ...
type JSONPathCompiled = *jp.Expr

// JSONPathParse ...
func JSONPathParse(jpath string) (JSONPathCompiled, error) {

	jpc, err := jp.ParseString(jpath)
	return &jpc, err
}
