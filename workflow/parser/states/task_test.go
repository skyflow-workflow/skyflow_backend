package states

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestTaskGetBone(t *testing.T) {

	var testcases = []struct {
		task    *Task
		except  StateBone
		wantErr bool
	}{
		{
			task: &Task{
				BaseState: &BaseState{
					Type: "Task",
					Next: "NextState",
				},
				TaskBody: &TaskBody{
					Catch: []TaskCatchNode{
						{
							Next: "CatchNextState",
						},
						{
							Next: "CatchNextState",
						},
					},
				},
			},
			except: StateBone{
				BaseBone: BaseBone{
					Type: "Task",
					Next: []string{"NextState", "CatchNextState", "CatchNextState"},
				},
			},
			wantErr: false,
		},
		{
			task: &Task{
				BaseState: &BaseState{
					Type: "Task",
					Next: "NextState",
				},
				TaskBody: &TaskBody{
					Catch: []TaskCatchNode{
						{
							Next: "CatchNextState",
						},
						{
							Next: "CatchNextState",
						},
					},
				},
			},
			except: StateBone{
				BaseBone: BaseBone{
					Type: "Task",
					Next: []string{"NextState", "CatchNextState", "CatchNextState"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range testcases {
		actualbone := tt.task.GetBone()
		assert.Equal(t, actualbone, tt.except)
	}

}

func TestTaskGetTaskTimeout(t *testing.T) {

}

func TestTaskGetOutput(t *testing.T) {

}

func TestTaskHasIntersection(t *testing.T) {

}

func TestTaskGetOutputWithPath(t *testing.T) {

}

func TestHasIntersection(t *testing.T) {

	var testcases = []struct {
		a      []string
		b      []string
		except bool
	}{
		{
			a:      []string{"a", "b", "c"},
			b:      []string{"a", "b", "c"},
			except: true,
		},
		{
			a:      []string{"a", "b", "c"},
			b:      []string{"a", "b", "c", "d"},
			except: true,
		},
		{
			a:      []string{"a", "b", "c"},
			b:      []string{"a", "b", "c", "d", "e"},
			except: true,
		},
		{
			a:      []string{"a", "b", "c"},
			b:      []string{"d", "e", "f"},
			except: false,
		},
	}

	for _, tt := range testcases {
		actual := HasIntersection(tt.a, tt.b)
		assert.Equal(t, actual, tt.except)
	}

}
