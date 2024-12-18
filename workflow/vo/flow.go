package vo

import (
	"github.com/mmtbak/microlibrary/paging"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
)

// ParseFlowRequest request
type ParseFlowRequest struct {
	WorkflowDefinition string
}

type ListNamespacesRequest struct {
	PageRequest paging.PageRequest
}

type ListNamespacesResponse struct {
	Namespaces   []po.Namespace
	PageResponse paging.PageResponse
}

type CreateActivityRequest struct {
	ActivityName string
	Comment      string
	Namespace    string
}

type CreateWorkflowRequest struct {
	WorkflowName string
	Comment      string
	Namespace    string
	Definition   string
}

type ListActivitiesRequest struct {
	PageRequest paging.PageRequest
}

type ListActivitiesResponse struct {
	Activities   []po.Function
	PageResponse paging.PageResponse
}

type ListWofkflowsRequest struct {
	PageRequest paging.PageRequest
}

type ListWorkflowsResponse struct {
	Workflows    []po.StateMachine
	PageResponse paging.PageResponse
}

type DescribeStepResponse struct {
	ExecutionUUID string
	Step          po.State
}
