package states

import (
	"fmt"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestSucceedState(t *testing.T) {

	var testcases = []struct {
		template  string
		wantError bool
	}{
		{
			template: `{
				"Type":"Succeed"
				}`,
			wantError: false,
		},
		{
			template: `{
				"Type":"Succeed",
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
	}
}
