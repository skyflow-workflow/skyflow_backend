package states

import (
	"strings"

	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
)

var myJson jsoniter.API

// myValidate self define validator
var myValidate = validator.New()

func init() {
	myJson = jsoniter.Config{
		EscapeHTML:    true,
		CaseSensitive: true, // 配置大小写敏感
	}.Froze()
}

// StringToMap 转换成map
func StringToMap(s string) (map[string]interface{}, error) {

	var v map[string]interface{}

	err := myJson.NewDecoder(strings.NewReader(s)).Decode(&v)
	return v, err

}

// Encoder 转码成 bytes
func Encoder(v interface{}) ([]byte, error) {
	return myJson.Marshal(v)
}

// ToString 转码成string 格式
func ToString(v interface{}) (string, error) {

	b, err := myJson.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), err
}

// MapUpdate merge map toi into map object i,
func MapUpdate(i map[string]interface{}, toi map[string]interface{}) {
	if i == nil {
		i = make(map[string]interface{})
	}
	for k, v := range toi {
		i[k] = v
	}
}
