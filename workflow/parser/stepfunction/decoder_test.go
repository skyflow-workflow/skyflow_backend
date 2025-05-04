package stepfunction

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestParserStateMachine(t *testing.T) {

	t.Log("Scanning examples directory")
	statemachinepath := "./examples/"
	direntries, err := os.ReadDir(statemachinepath)
	assert.Equal(t, err, nil)
	for _, direntry := range direntries {
		t.Run(fmt.Sprintf("Test_%s", direntry.Name()), func(t *testing.T) {
			if !direntry.Type().IsRegular() {
				t.Logf("'%s' is not regular file\n", direntry.Name())
				return
			}
			filename := direntry.Name()
			ext := filepath.Ext(direntry.Name())
			if ext != ".json" {
				t.Logf("'%s' is not a JSON file", direntry.Name())
				return
			}
			// TestParserStateMachine test the parsing of a state machine
			// from a JSON definition
			//
			// The definition is a JSON string that contains the state machine
			// definition
			//
			// The function should return a StateMachine object
			//
			// The function should return an error if the definition is invalid
			// or if the state machine is not valid
			//
			filepath := filepath.Join(statemachinepath, filename)
			t.Log("Parsing File:", filepath)
			filecontent, err := os.ReadFile(filepath)
			assert.Equal(t, err, nil)
			decoder := NewStepfuncionDecoder(nil, nil)
			_, err = decoder.Decode(string(filecontent))
			if err != nil {
				t.Logf("Error parsing file: %s", err)
				t.Fail()
				return
			}
			assert.Equal(t, err == nil, true)
		},
		)
	}
}
