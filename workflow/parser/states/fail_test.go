package states

import (
	"fmt"
	"strings"
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

// TestFailState_GetDefinition_NoDeniedFields ensures GetDefinition does not include OutputPath/InputPath
// so that the stored step definition can be parsed later (e.g. in DescribeExecutionBone) without "field is denied".
func TestFailState_GetDefinition_NoDeniedFields(t *testing.T) {
	state := &FailState{
		BaseState: &BaseState{
			Name:       "MyFail",
			Type:       "Fail",
			OutputPath: "$",
			InputPath:  "$",
			End:        true,
		},
		FailBody: &FailBody{
			Error: "SomeError",
			Cause: "SomeCause",
			Abort: false,
		},
	}
	def, err := state.GetDefinition()
	assert.Equal(t, err, nil)
	assert.Equal(t, strings.Contains(def, "OutputPath"), false)
	assert.Equal(t, strings.Contains(def, "InputPath"), false)
	// Def must be flat and contain Type so round-trip works
	assert.Equal(t, strings.Contains(def, "Type"), true)
	data, err := StringToMap(def)
	assert.Equal(t, err, nil)
	_, hasType := data["Type"]
	assert.Equal(t, hasType, true)
	// Round-trip: parsed definition must not trigger "field is denied"
	_, err = NewFailStateFromString(def)
	assert.Equal(t, err, nil)
}
