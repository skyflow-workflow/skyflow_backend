package executor

import (
	"fmt"
	"testing"
)

func TestMergeExecutionStatus(t *testing.T) {

	statuses := []string{"Created",
		// "WaitInit",
		"Failed",
		"Failed",
		// "Running",
	}

	status := MergeExecutionStatus(statuses)
	fmt.Println("finial status:  ", status)

}
