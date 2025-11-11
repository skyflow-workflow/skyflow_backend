package dispatcher

import (
	"fmt"
	"testing"

	"github.com/skyflow-workflow/skyflow_backbend/workflow"
)

func TestCreateDispatcherService(t *testing.T) {

	var wfSvc workflow.WorkflowService

	dpSvc, err := NewDispatcher(wfSvc, DefaultOption)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(dpSvc)

}
