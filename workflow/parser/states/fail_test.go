package states

import (
	"fmt"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestRunFailState(t *testing.T) {

	var testcases = []struct {
		state    *FailState
		faildata FailData
	}{
		{
			state: &FailState{
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
				Error: "test error1",
			},
		},
		{
			state: &FailState{
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

func TestParseFailState(t *testing.T) {

	var testcases = []struct {
		defintion string
		wantState *FailState
		wantError bool
	}{
		{
			defintion: `{
				"Type":"Fail"
				}`,
			wantState: &FailState{
				BaseState: &BaseState{
					Type:            "Fail",
					OutputPath:      "$",
					MaxExecuteTimes: 1000,
					End:             true,
					Next:            "",
				},
				FailBody: &FailBody{
					Abort: false,
					Cause: "",
					Error: "",
				},
			},
			wantError: false,
		},
		{
			defintion: `{
				"Type":"Fail",
				"Error": "这个错误我处理不了",
				"Cause": "这个代码bug了"
				}`,
			wantState: &FailState{
				BaseState: &BaseState{
					Type:            "Fail",
					Next:            "",
					End:             true,
					OutputPath:      "$",
					MaxExecuteTimes: 1000,
				},
				FailBody: &FailBody{
					Abort: false,
					Cause: "这个代码bug了",
					Error: "这个错误我处理不了",
				},
			},
			wantError: false,
		},
		{
			defintion: `{
				"Type":"Fail",
				"Next":""
				}`,
			wantState: &FailState{
				BaseState: &BaseState{
					Type:            "Fail",
					Next:            "X",
					End:             true,
					OutputPath:      "$",
					MaxExecuteTimes: 1000,
				},
				FailBody: &FailBody{
					Abort: false,
					Cause: "",
					Error: "",
				},
			},
			wantError: true,
		},
	}

	for idx, tt := range testcases {
		t.Run(fmt.Sprintf("index - %d", idx), func(t *testing.T) {
			state, err := NewFailStateFromString(tt.defintion)
			assert.Equal(t, err != nil, tt.wantError)
			if err != nil {
				return
			}
			assert.Equal(t, state.GetBaseState(), tt.wantState.BaseState)
			assert.Equal(t, state.FailBody, tt.wantState.FailBody)
		})
	}
}
