package states

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestTaskGetNextState(t *testing.T) {

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
