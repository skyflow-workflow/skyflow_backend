package stepfunction

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestParserStateMachine(t *testing.T) {

	println("Scanning examples directory")
	statemachinepath := "./examples/"
	direntries, err := os.ReadDir(statemachinepath)
	assert.Equal(t, err, nil)
	for _, direntry := range direntries {
		if !direntry.Type().IsRegular() {
			println(direntry.Name(), " is not regular file ,Pass")
			continue
		}
		filename := direntry.Name()
		ext := filepath.Ext(direntry.Name())
		if ext != ".json" {
			println(direntry.Name(), " is not a JSON file, Pass ")
			continue
		}
		t.Run(fmt.Sprintf("Test_%s", filename), func(t *testing.T) {
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
				t.Log(err)
			}
			assert.Equal(t, err == nil, true)
		},
		)
	}
}
