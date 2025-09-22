package states

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
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

		t.Run(fmt.Sprintf("index - %d", idx), func(t *testing.T) {
			_, err := NewSucceedStateFromString(tt.template)
			assert.Equal(t, tt.wantError, err != nil)

		})
	}
}
