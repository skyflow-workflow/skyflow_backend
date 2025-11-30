package apiserver

import (
	"github.com/skyflow-workflow/skyflow_backend/workflow"
)

// APIServer is the API server.
type APIServer struct {
	workflowSvc workflow.WorkflowService
}

// NewAPIServer creates a new API server.
func NewAPIServer(svc workflow.WorkflowService) *APIServer {

	return &APIServer{
		workflowSvc: svc,
	}
}
