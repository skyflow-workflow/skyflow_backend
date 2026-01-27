package executor

import (
	"log/slog"
	"testing"

	"github.com/skyflow-workflow/skyflow_backend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"gopkg.in/go-playground/assert.v1"
)

func Test_BuildExecutionBone(t *testing.T) {

	var dbSteps = []po.Step{
		{
			ID:          1,
			ExecutionID: 1,
			Type:        "Task",
			Name:        "Step1",
			Status:      "Running",
			Definition: `{
				"Type": "Task",
				"Resource": "activity:function1",
				"Next": "Step2",
				"Comment": "test comment for step1"
			}`,
			GroupIndex:   1,
			GroupID:      1,
			ExecuteIndex: 11,
		},
		{
			ID:          2,
			ExecutionID: 1,
			Type:        "Task",
			Name:        "Step2",
			Status:      "Running",
			Definition: `{
				"Type": "Task",
				"Resource": "activity:function1",
				"Next": "Step2",
				"Comment": "test comment for step2"
			}`,
			GroupIndex:   1,
			GroupID:      1,
			ExecuteIndex: -1,
		},
	}

	req := BuildExecutionBoneRequest{
		Steps: &dbSteps,
	}

	exeBone, err := BuildExecutionBone(req, StandardExecutor)
	assert.Equal(t, err, nil)
	slog.Info("Execution Bone", "exeBone", exeBone)
	boneStr, err := toolkit.ToString(exeBone)
	assert.Equal(t, err, nil)
	slog.Info("Execution Bone", "exeBoneJSON", boneStr)
}
