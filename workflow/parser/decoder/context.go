package decoder

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

// DecoderPath ...
type DecoderPath struct{}

// AddPath ...
func AddPath(ctx context.Context, paths ...string) context.Context {
	curpath := GetPath(ctx)
	curpath = append(curpath, paths...)
	ctx = context.WithValue(ctx, DecoderPath{}, curpath)
	return ctx
}

// GetPath ...
func GetPath(ctx context.Context) []string {
	defautpath := []string{}
	if ctx == nil {
		return defautpath
	}
	if val := ctx.Value(DecoderPath{}); val != nil {
		return val.([]string)
	}
	return defautpath
}

// MergeError ...
func MergeError(ctx context.Context, err error) error {

	path := GetPath(ctx)

	if err, ok := err.(*states.FieldError); ok {
		newerr := &states.FieldError{
			RawError: err,
			Paths:    append(path, err.Paths...),
		}
		return newerr
	}
	return err
}
