package states

import (
	"fmt"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestToMap(t *testing.T) {

	var testcases = []struct {
		template  string
		wantError bool
	}{
		{
			template:  `{"ab": 10 }`,
			wantError: false,
		}, {
			template:  `{"x": x }`,
			wantError: true,
		},
	}
	for idx, tt := range testcases {
		fmt.Println("index :", idx)
		val, err := StringToMap(tt.template)
		fmt.Println(err)
		fmt.Println(val)
		assert.Equal(t, err != nil, tt.wantError)
	}
}

func TestMapUpdate(t *testing.T) {

	i := map[string]interface{}{
		"a": 10,
		"b": 20,
	}
	toi := map[string]interface{}{
		"a": 100,
		"c": 30,
	}
	MapUpdate(i, toi)
	assert.Equal(t, i["a"], 100)
	assert.Equal(t, i["b"], 20)
	assert.Equal(t, i["c"], 30)

}
