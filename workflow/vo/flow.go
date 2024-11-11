package vo

import (
	"gopkg.mihoyo.com/plat/cloudflow/pkg/paging"
	"gopkg.mihoyo.com/plat/cloudflow/workflow/po"
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
	Activities   []po.Activity
	PageResponse paging.PageResponse
}

type ListWofkflowsRequest struct {
	PageRequest paging.PageRequest
}

type ListWorkflowsResponse struct {
	Workflows    []po.Workflow
	PageResponse paging.PageResponse
}

type DescribeStepResponse struct {
	ExecutionUUID string
	Step          po.Step
}
