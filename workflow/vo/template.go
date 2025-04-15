package vo

import (
	"github.com/mmtbak/microlibrary/paging"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
)

// ParseFlowRequest request
type ParseFlowRequest struct {
	StateMachineDefinition string
}

type CreateNamespaceRequest struct {
	Name    string
	Comment string
}

type CreateNamespaceResponse struct {
	Data po.Namespace
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

type DeleteNamespaceRequest struct {
	Name string
}

// CreateActivityRequest ...
type CreateActivityRequest struct {
	ActivityName string
	Comment      string
	Namespace    string
}

type CreateActivityResponse struct {
	Data po.Activity
}

// CreateStateMachineRequest ...
type CreateStateMachineRequest struct {
	StateMachineName string
	Comment          string
	Namespace        string
	Definition       string
}

type CreateStateMachineResponse struct {
	Data po.StateMachine
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

// DescribeActivityRequest ...
type DescribeActivityRequest struct {
	ActivityURI string
}

type DescribeActivityResponse struct {
	ActivityURI string
	Name        string
	Comment     string
	CreateTime  int64
	UpdateTime  int64
}

type DeleteActivityRequest struct {
	ActivityURI string
}

// ListStateMachinesRequest ...
type ListStateMachinesRequest struct {
	PageRequest paging.PageRequest
}

// ListStateMachinesResponse ...
type ListStateMachinesResponse struct {
	StateMachines []po.StateMachine
	PageResponse  paging.PageResponse
}

// DescribeStepResponse ...
type DescribeStepResponse struct {
	ExecutionUUID string
	Step          po.State
}
