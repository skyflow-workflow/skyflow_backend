package vo

import (
	"github.com/mmtbak/microlibrary/paging"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
)

// ParseFlowRequest request
type ParseFlowRequest struct {
	WorkflowDefinition string
}

type CreateNamespaceRequest struct {
	Name    string
	Comment string
}

type CreateNamespaceResponse struct {
	po.Namespace
}

// ListNamespacesRequest ...
type ListNamespacesRequest struct {
	PageRequest paging.PageRequest
}

// ListNamespacesResponse ...
type ListNamespacesResponse struct {
	Namespaces   []po.Namespace
	PageResponse paging.PageResponse
}

// CreateActivityRequest ...
type CreateActivityRequest struct {
	ActivityName string
	Comment      string
	Namespace    string
}

type CreateActivityResponse struct {
	po.Activity
}

// CreateWorkflowRequest ...
type CreateWorkflowRequest struct {
	WorkflowName string
	Comment      string
	Namespace    string
	Definition   string
}

// ListActivitiesRequest ...
type ListActivitiesRequest struct {
	PageRequest paging.PageRequest
}

// ListActivitiesResponse ...
type ListActivitiesResponse struct {
	Activities   []po.Activity
	PageResponse paging.PageResponse
}

// ListWofkflowsRequest ...
type ListWofkflowsRequest struct {
	PageRequest paging.PageRequest
}

// ListWorkflowsResponse ...
type ListWorkflowsResponse struct {
	Workflows    []po.StateMachine
	PageResponse paging.PageResponse
}

// DescribeStepResponse ...
type DescribeStepResponse struct {
	ExecutionUUID string
	Step          po.State
}
