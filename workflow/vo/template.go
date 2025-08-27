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
	Name        string
	Description string
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
	Description  string
	Namespace    string
	// Parameters is the parameters description of the activity, it is a json string
	Parameters string
}

type CreateActivityResponse struct {
	Data po.Activity
}

type CreateStateMachineRequest struct {
	Name        string
	Description string
	Namespace   string
	Definition  string
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
	Description string
	CreateTime  int64
	UpdateTime  int64
}

type DeleteActivityRequest struct {
	ActivityURI string
}

// ListStateMachinesRequest ...
type ListStateMachinesRequest struct {
	Namespace   string
	PageRequest paging.PageRequest
}

// ListStateMachinesResponse ...
type ListStateMachinesResponse struct {
	StateMachines []po.StateMachine
	PageResponse  paging.PageResponse
}

type DescribeStateMachineRequest struct {
	StateMachineURI string
}
type DescribeStateMachineResponse struct {
	Data *po.StateMachine
}

type DeleteStateMachineRequest struct {
	StateMachineURI string
}

// DescribeStepResponse ...
type DescribeStepResponse struct {
	ExecutionUUID string
	Step          po.Step
}

type UpdateStateMachineRequest struct {
	StateMachineURI string
	Name            string
	Description     string
	Definition      string
}
