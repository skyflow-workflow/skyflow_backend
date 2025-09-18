package states

import (
	"fmt"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestParseFailState(t *testing.T) {

	var testcases = []struct {
		template  string
		wantError bool
	}{
		{
			template: `{
				"Type":"Fail"
				}`,
			wantError: false,
		},
		{
			template: `{
				"Type":"Fail",
				"Error": "这个错误我处理不了",
				"Cause": "这个代码bug了"
				}`,
			wantError: false,
		},
		{
			template: `{
				"Type":"Fail",
				"Next":"X"
				}`,
			wantError: true,
		},
	}

	for idx, tt := range testcases {

		fmt.Println("index :", idx)
		state, err := NewFailStateFromString(tt.template)
		fmt.Println(err)
		fmt.Println(state)
		assert.Equal(t, tt.wantError, err != nil)
		if err == nil {
			data := state.GetFailData()
			fmt.Println(data)
			bone := state.GetBone()
			fmt.Println(bone)
		}
	}
}
