package decoder

import (
	"context"
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"

	jsoniter "github.com/json-iterator/go"
)

var myJson jsoniter.API

func init() {
	myJson = jsoniter.Config{
		EscapeHTML:             true,
		CaseSensitive:          true, // 配置大小写敏感
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
	}.Froze()
}

// CommonDecoder ...
type CommonDecoder struct {
}

// NewCommonDecoder ...
func NewCommonDecoder() *CommonDecoder {
	return &CommonDecoder{}
}

// Decode ...
func (decoder *CommonDecoder) Decode(definition string) (*states.StateMachine, error) {
	return nil, nil
}

// JSONUnmarshall unmarshal the json string to the target object
func (decoder *CommonDecoder) JSONUnmarshall(data string, v any) error {
	err := myJson.Unmarshal([]byte(data), v)
	if err != nil {
		if jsonerr, ok := err.(*json.SyntaxError); ok {
			newerr := &states.FieldError{
				RawError: err,
				Offset:   jsonerr.Offset,
			}
			return newerr
		}
		return err
	}
	return nil
}

// MapDecode decode the map to the target object
func (decoder *CommonDecoder) MapDecode(input any, output any) error {
	err := mapstructure.Decode(input, output)
	return err
}

// GetCtxDecodePath ...
func (decoder *CommonDecoder) GetCtxDecodePath(ctx context.Context) []string {
	return GetPath(ctx)
}

// AddCtxDecodePath ...
func (decoder *CommonDecoder) AddCtxDecodePath(ctx context.Context, paths ...string) context.Context {
	return AddPath(ctx, paths...)
}

// NewFieldPathError ...
func (decoder *CommonDecoder) NewFieldPathError(ctx context.Context, err error) error {
	return NewFieldPathError(ctx, err)
}
