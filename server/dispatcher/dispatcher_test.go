package dispatcher

import (
	"fmt"
	"testing"

	"github.com/skyflow-workflow/skyflow_backbend/workflow"
)

func TestCreateDispatcherService(t *testing.T) {

	var wfsvc workflow.WorkflowService

	dpsvc, err := NewDispatcher(wfsvc, DefaultOption)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(dpsvc)

}
