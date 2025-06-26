package workflow

import (
	"github.com/mmtbak/microlibrary/rdb"
)

// WorkflowService Service provides workflow related services
// It is used to manage workflow execution, templates, and events.
// It is also used to manage workflow state and activity tasks.
// The WorkflowService interface defines the methods that the service should implement.

type WorkflowService *workflowService

// workflowService is the implementation of WorkflowService interface
// It provides methods to manage workflow execution, templates, and events.
type workflowService struct {
	DBClient *rdb.DBClient
}

func NewWorkflowService(dbClient *rdb.DBClient) WorkflowService {
	return &workflowService{}
}
