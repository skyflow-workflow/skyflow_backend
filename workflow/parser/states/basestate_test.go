package states

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestValidateStateFieldOptional(t *testing.T) {

	tests := []struct {
		name string
		data map[string]any
		want error
	}{
		{
			name: "lack of next or end",
			data: map[string]any{
				"Type":    "Task",
				"Comment": "comment",
			},
			want: NewFieldPathError(fmt.Errorf("%w: Next or End", ErrorLackOfRequiredField), StateFieldNames.Next),
		},
		{
			name: "next and end should not be both defined",
			data: map[string]any{
				"Type":       "Task",
				"Comment":    "comment",
				"InputPath":  "inputpath",
				"OutputPath": "outputpath",
				"Parameters": "parameters",
				"ResultPath": "resultpath",
				"Next":       "next",
				"End":        "end",
			},
			want: NewFieldPathError(fmt.Errorf("%w: Next and End should not be both defined", ErrorInvalidField), StateFieldNames.Next),
		},
		{
			name: "choice next end is deny",
			data: map[string]any{
				"Type":    "Choice",
				"Comment": "comment",
				"Next":    "next",
			},
			want: NewFieldPathError(fmt.Errorf("%w: Next", ErrorFiledDenied), StateFieldNames.Next),
		},
		{
			name: "choice end is deny",
			data: map[string]any{
				"Type":    "Choice",
				"Comment": "comment",
				"End":     "end",
			},
			want: NewFieldPathError(fmt.Errorf("%w: End", ErrorFiledDenied), StateFieldNames.End),
		},
		{
			name: "success next is deny",
			data: map[string]any{
				"Type":    "Succeed",
				"Comment": "comment",
				"Next":    "next",
			},
			want: NewFieldPathError(fmt.Errorf("%w: Next", ErrorFiledDenied), StateFieldNames.Next),
		},
		{
			name: "success end is deny",
			data: map[string]any{
				"Type":    "Succeed",
				"Comment": "comment",
				"End":     "end",
			},
			want: NewFieldPathError(fmt.Errorf("%w: End", ErrorFiledDenied), StateFieldNames.End),
		},
		{
			name: "pass next is optional",
			data: map[string]any{
				"Type":    "Pass",
				"Comment": "comment",
				"Next":    "next",
			},
			want: nil,
		},
		{
			name: "pass end is optional",
			data: map[string]any{
				"Type":    "Pass",
				"Comment": "comment",
				"End":     "end",
			},
			want: nil,
		},
		{
			name: "pass next/end is deny",
			data: map[string]any{
				"Type":    "Pass",
				"Comment": "comment",
				"Next":    "next",
				"End":     "end",
			},
			want: NewFieldPathError(fmt.Errorf("%w: Next and End should not be both defined", ErrorInvalidField), StateFieldNames.Next),
		},
		{
			name: "wait next is optional",
			data: map[string]any{
				"Type":    "Wait",
				"Comment": "comment",
				"Next":    "next",
			},
			want: nil,
		},
		{
			name: "wait end is optional",
			data: map[string]any{
				"Type":    "Wait",
				"Comment": "comment",
				"End":     "end",
			},
			want: nil,
		},
		{
			name: "wait next/end is deny",
			data: map[string]any{
				"Type":    "Wait",
				"Comment": "comment",
				"Next":    "next",
				"End":     "end",
			},
			want: NewFieldPathError(fmt.Errorf("%w: Next and End should not be both defined", ErrorInvalidField), StateFieldNames.Next),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStateFieldOptional(test.data)
			if err == nil {
				assert.Equal(t, test.want == nil, true)
			} else {
				if test.want == nil {
					t.Log(err.Error())
				}
				assert.Equal(t, err.Error(), test.want.Error())
			}
		})
	}
}

func TestBaseStateInit(t *testing.T) {
	tests := []struct {
		name      string
		baseState *BaseState
		want      error
	}{
		{
			name: "task compile failed",
			baseState: &BaseState{
				Type:       "Task",
				Comment:    "comment",
				Next:       "next",
				InputPath:  "input",
				OutputPath: "output",
				ResultPath: "$.result",
				Parameters: "$.parameters",
				Retry:      "$.retry",
				Catch:      "$.catch",
			},
			want: nil,
		},
		{
			name: "task compile success",
			baseState: &BaseState{
				Type:       "Task",
				Comment:    "comment",
				Next:       "next",
				InputPath:  "$.input",
				OutputPath: "$.output",
				ResultPath: "$.result",
				Parameters: "$.parameters",
				Retry:      "$.retry",
				Catch:      "$.catch",
			},
			want: nil,
		},
		{
			name: "choice compile",
			baseState: &BaseState{
				Type:       "Choice",
				Comment:    "comment",
				InputPath:  "$.input",
				OutputPath: "$.output",
				ResultPath: "$.result",
				Parameters: "$.parameters",
				Retry:      "$.retry",
				Catch:      "$.catch",
			},
			want: nil,
		},
		{
			name: "pass",
			baseState: &BaseState{
				Type:    "Pass",
				Comment: "comment",
				Next:    "next",
			},
			want: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.baseState.Init()
			assert.Equal(t, got, test.want)
		})
	}
}

func TestBaseStateGetBone(t *testing.T) {
	tests := []struct {
		name      string
		baseState *BaseState
		want      StateBone
	}{
		{
			name: "task",
			baseState: &BaseState{
				Type:       "Task",
				Comment:    "comment",
				Next:       "next",
				End:        false,
				InputPath:  "$.input",
				OutputPath: "$.output",
				ResultPath: "$.result",
			},
			want: StateBone{
				BaseBone: BaseBone{
					Type:    "Task",
					Comment: "comment",
					Next:    []string{"next"},
				},
			},
		},
		{
			name: "choice",
			baseState: &BaseState{
				Type:       "Choice",
				Comment:    "Choice State Comment",
				Next:       "",
				End:        false,
				InputPath:  "$.input",
				OutputPath: "$.output",
				ResultPath: "$.result",
			},
			want: StateBone{
				BaseBone: BaseBone{
					Type:    "Choice",
					Comment: "Choice State Comment",
					Next:    []string{},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.baseState.GetBone()
			assert.Equal(t, got, test.want)
		})
	}
}

func TestGetParametersInput(t *testing.T) {

	tests := []struct {
		name  string
		state *BaseState
		input any
		want  any
	}{
		{
			name: "parameters is nil",
			state: &BaseState{
				Type:       "Task",
				Parameters: nil,
			},
			input: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			want: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "parameters is not nil",
			state: &BaseState{
				Type: "Task",
				Parameters: map[string]any{
					"a.$": "$.key1",
					"b.$": "$.key2",
				},
			},
			input: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			want: map[string]any{
				"a": "value1",
				"b": "value2",
			},
		},

		{
			name: "inputpath is not nil",
			state: &BaseState{
				Type:      "Task",
				InputPath: "$.input",
				Parameters: map[string]any{
					"a.$": "$.key1",
					"b.$": "$.key2",
				},
			},
			input: map[string]any{
				"input": map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
			},
			want: map[string]any{
				"a": "value1",
				"b": "value2",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.state.GetParametersInput(test.input)
			if err != nil {
				t.Errorf("GetParametersInput() error = %v", err)
				return
			}
			assert.Equal(t, got, test.want)
		})
	}
}

func TestGetOutput(t *testing.T) {

	tests := []struct {
		name       string
		state      *BaseState
		input      any
		taskOutput any
		want       any
	}{
		{
			name: "outputpath is nil",
			state: &BaseState{
				Type:       "Task",
				OutputPath: "",
			},
			input: map[string]any{
				"output": map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
			},
			taskOutput: "newoutputvalue",
			want:       nil,
		},
		{
			name: "outputpath is not nil",
			state: &BaseState{
				Type:       "Task",
				OutputPath: "$.output",
			},
			input: map[string]any{
				"output": map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
			},
			taskOutput: "newoutputvalue",
			want: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "resultpath is not nil",
			state: &BaseState{
				Type:       "Task",
				ResultPath: "$.result",
				OutputPath: "$",
			},
			input: map[string]any{
				"output": map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
			},
			taskOutput: "taskresultvalue",
			want: map[string]any{
				"output": map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
				"result": "taskresultvalue",
			},
		},
		{
			name: "resultpath/outputpath is not nil",
			state: &BaseState{
				Type:       "Task",
				ResultPath: "$.result",
				OutputPath: "$.result",
			},
			input: map[string]any{
				"output": map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
			},
			taskOutput: "taskresultvalue",
			want:       "taskresultvalue",
		},
		{
			name: "state output map",
			state: &BaseState{
				Type:       "Task",
				ResultPath: "$.result",
				OutputPath: "$.result",
			},
			input: map[string]any{
				"output": map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
			},
			taskOutput: map[string]any{
				"key3": "value3",
				"key4": "value4",
			},
			want: map[string]any{
				"key3": "value3",
				"key4": "value4",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.state.GetOutput(test.input, test.taskOutput)
			assert.Equal(t, err, nil)
			assert.Equal(t, result, test.want)
		})
	}
}
