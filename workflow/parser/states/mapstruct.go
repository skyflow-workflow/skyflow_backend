package states

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

var myMapdecodeconfig = mapstructure.DecoderConfig{
	Metadata:             nil,
	IgnoreUntaggedFields: true,
	TagName:              "mapstructure",
}

var myMapDecoder *mapstructure.Decoder

func DecodeMapToStruct(input map[string]interface{}, output interface{}) error {

	config := myMapdecodeconfig
	config.Result = output

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func DecodeStructToMap(input interface{}) (map[string]interface{}, error) {

	var err error
	if input == nil {
		return nil, fmt.Errorf("input should not nil")
	}
	var output = map[string]interface{}{}
	config := myMapdecodeconfig
	config.Result = &output

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(input)
	if err != nil {
		return nil, err
	}

	return output, nil
}
