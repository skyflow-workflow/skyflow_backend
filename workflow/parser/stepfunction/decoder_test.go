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
			print("'%s' is not regular file ,Pass ", direntry.Name())
			continue
		}
		filename := direntry.Name()
		ext := filepath.Ext(direntry.Name())
		if ext != ".json" {
			print("'%s' is not a JSON file, Pass ", direntry.Name())
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
			println("Parsing file ", filepath)
			filecontent, err := os.ReadFile(filepath)
			assert.Equal(t, err, nil)
			decoder := NewStepfuncionDecoder(nil, nil)
			_, err = decoder.Decode(string(filecontent))
			assert.Equal(t, err, nil)
		},
		)
	}
}
