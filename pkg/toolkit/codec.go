package toolkit

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var myJson jsoniter.API

func init() {
	myJson = jsoniter.Config{
		EscapeHTML:    true,
		CaseSensitive: true, // 配置大小写敏感
	}.Froze()
}

// Encode encode data to json serialized bytes
func Encoder(v any) ([]byte, error) {
	return myJson.Marshal(v)
}

func Decode(data []byte, v any) error {
	return myJson.Unmarshal(data, v)
}

func DecodeString(data string, v any) error {
	return myJson.UnmarshalFromString(data, v)
}

// EncodeToString encode data to json serialized string
func EncodeToString(v any) (string, error) {

	b, err := myJson.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), err
}

// DecodeDataToMap decode data to map type
func DecodeStringToMap(s string) (map[string]any, error) {

	var v map[string]any
	err := myJson.NewDecoder(strings.NewReader(s)).Decode(&v)
	return v, err

}

func ToString(v any) (string, error) {
	b, err := myJson.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
