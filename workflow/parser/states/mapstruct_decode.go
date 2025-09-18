package states

import "github.com/mitchellh/mapstructure"

var mymapdecodeconfig = mapstructure.DecoderConfig{
	Metadata:             nil,
	IgnoreUntaggedFields: true,
}

func MapStructDecode(input interface{}, output interface{}) error {

	config := mymapdecodeconfig
	config.Result = output

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}
