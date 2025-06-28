package executor

import (
	"fmt"
	"testing"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"

	"github.com/go-playground/assert/v2"
)

func init() {

}
func TestInitState(t *testing.T) {
	step_id := 2541

	var dbStep *po.Step
	var err error
	myExecutor := StandardExecutor

	dbStep, err = myExecutor.QueryStepByID(step_id, []string{}, nil)

	assert.Equal(t, err, nil)
	es, err := NewExecutionStep(dbStep, myExecutor)
	fmt.Println(err)
	assert.Equal(t, err, nil)
	es.GetInput()
}
