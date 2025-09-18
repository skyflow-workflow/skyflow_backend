package states

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestParseFailState(t *testing.T) {

	var testcases = []struct {
		state    *Fail
		faildata FailData
	}{
		{
			state: &Fail{
				BaseState: &BaseState{
					Name: "fail1",
					Type: "Fail",
					Next: "end",
				},
				FailBody: &FailBody{
					Abort: true,
					Cause: "test cause1",
					Error: "test error1",
				},
			},
			faildata: FailData{
				Cause: "test cause1",
				Error: "test error2",
			},
		},
		{
			state: &Fail{
				BaseState: &BaseState{
					Name: "fail2",
					Type: "Fail",
					Next: "end",
				},
				FailBody: &FailBody{
					Abort: true,
					Cause: "test cause2",
					Error: "test error2",
				},
			},
			faildata: FailData{
				Cause: "test cause2",
				Error: "test error2",
			},
		},
	}

	for _, tt := range testcases {
		assert.Equal(t, tt.state.GetFailData(), tt.faildata)
	}
}
