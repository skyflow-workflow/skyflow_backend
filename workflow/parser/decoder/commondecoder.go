package decoder

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"

	jsoniter "github.com/json-iterator/go"
)

var myjson jsoniter.API

func init() {
	myjson = jsoniter.Config{
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

// JSONUnmashall ...
func (decoder *CommonDecoder) JSONUnmashall(data string, v any) error {
	err := myjson.Unmarshal([]byte(data), v)
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

// MapDecode ...
func (decoder *CommonDecoder) MapDecode(input interface{}, output interface{}) error {
	err := mapstructure.Decode(input, output)
	return err
}
