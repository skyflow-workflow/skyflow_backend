// package template
package template

import (
	"context"
	"fmt"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// TemplateService template management service
type TemplateService = *templateService

type templateService struct {
	metadb *rdb.DBClient
}

func NewTemplateService(metadb *rdb.DBClient) TemplateService {
	svc := &templateService{
		metadb: metadb,
	}
	return svc
}

func (svc *templateService) CreateNamespace(ctx context.Context, req vo.CreateNamespaceRequest, tx rdb.Tx) (vo.CreateNamespaceResponse, error) {

	var ns po.Namespace
	var err error

	tx, maker := svc.metadb.NewTxMaker(tx)
	defer maker.Close(&err)

	err = tx.Where("name = ?", req.Name).First(&ns).Error
	if err != nil && !rdb.IsRecordNotFound(err) {
		return vo.CreateNamespaceResponse{}, err
	}

	newns := po.Namespace{
		Name:    name,
		Comment: comment,
	}
	// 如果存在则更新，否则创建,
	// 为什么不用Save， Save会更新所有字段，而我们只需要更新部分字段
	if ns.ID != 0 {
		err = tx.Where(po.Namespace{ID: ns.ID}).Updates(&newns).Error
	} else {
		err = tx.Create(&newns).Error
	}

	return newns, err
}

// ListNamespaces implements cloudflow.CloudflowServer.
func (svc *templateService) ListNamespaces(ctx context.Context, req vo.ListNamespacesRequest) (vo.ListNamespacesResponse, error) {

	var err error
	var count int64
	var nss = []po.Namespace{}
	var resp vo.ListNamespacesResponse
	tx, maker := svc.metadb.NewTxMaker(nil)
	defer maker.Close(&err)
	limit, offset := req.PageRequest.Limit()

	tx = tx.Model(new(po.Namespace))
	err = tx.Count(&count).Error
	if err != nil {
		return resp, err
	}
	err = tx.Order("name asc").Offset(offset).Limit(limit).Find(&nss).Error
	if err != nil {
		return resp, err
	}
	pageresp := req.PageRequest.Response(count)
	resp.Namespaces = nss
	resp.PageResponse = pageresp

	return resp, err
}

// CreateActivity implements cloudflow.CloudflowServer.
func (svc *templateService) CreateActivity(ctx context.Context, req vo.CreateActivityRequest, tx rdb.Tx) (po.Activity, error) {

	var ns po.Namespace
	var activity po.Activity
	var err error

	tx, maker := svc.metadb.NewTxMaker(tx)
	defer maker.Close(&err)

	activity_uri := fmt.Sprintf("%s:%s/%s", "activity", req.Namespace, req.ActivityName)
	ns, err = svc.CreateNamespace(ctx, req.Namespace, "", tx)
	if err != nil {
		return activity, err
	}

	activitycond := po.Activity{
		URI: activity_uri,
	}

	err = tx.Model(new(po.Activity)).Where(activitycond).Select("id").First(&activity).Error
	if err != nil && !rdb.IsRecordNotFound(err) {
		return activity, err
	}

	newactivity := po.Activity{
		NamespaceID: ns.ID,
		Name:        req.ActivityName,
		Comment:     req.Comment,
		URI:         activity_uri,
		Status:      ActivityStatus.Enable,
	}

	if activity.ID != 0 {
		err = tx.Where(po.Activity{ID: activity.ID}).Updates(&newactivity).Error
	} else {
		err = tx.Create(&newactivity).Error
	}

	if err != nil {
		return activity, err
	}

	return newactivity, nil
}

// CreateWorkflow implements cloudflow.CloudflowServer.
func (svc *templateService) CreateWorkflow(ctx context.Context, req vo.CreateWorkflowRequest, tx rdb.Tx) (po.StateMachine, error) {

	var ns po.Namespace
	var workflow po.StateMachine
	var err error

	flow, err := parser.ParseStateMachine(req.Definition)
	if err != nil {
		return workflow, err
	}

	workflowuri := fmt.Sprintf("%s:%s/%s", "workflow", req.Namespace, req.WorkflowName)

	tx, maker := svc.metadb.NewTxMaker(tx)
	defer maker.Close(&err)

	ns, err = svc.CreateNamespace(ctx, req.Namespace, "", tx)
	if err != nil {
		return workflow, err
	}

	workfowcond := po.StateMachine{
		// NamespaceID: ns.ID,
		// Name:        req.WorkflowName,
		URI: workflowuri,
	}

	err = tx.Model(new(po.StateMachine)).Where(workfowcond).Select("id").First(&workflow).Error
	if err != nil && !rdb.IsRecordNotFound(err) {
		return workflow, err
	}

	newworkflow := po.StateMachine{
		NamespaceID: ns.ID,
		Name:        req.WorkflowName,
		Type:        flow.Type,
		Comment:     req.Comment,
		URI:         workflowuri,
		Definition:  req.Definition,
		Status:      ActivityStatus.Enable,
	}
	if workflow.ID != 0 {
		err = tx.Where(po.StateMachine{ID: workflow.ID}).Updates(&newworkflow).Error
	} else {
		err = tx.Create(&newworkflow).Error
	}
	if err != nil {
		return workflow, err
	}
	return newworkflow, nil
}

// DescribeActivity implements cloudflow.CloudflowServer.
func (svc *templateService) DescribeActivity(ctx context.Context, activityuri string, tx rdb.Tx) (po.Activity, error) {
	var activity po.Activity
	var err error
	tx, maker := svc.metadb.NewTxMaker(tx)
	defer maker.Close(&err)

	err = tx.Where(po.Activity{URI: activityuri}).Take(&activity).Error
	return activity, err
}

// DescribeWorkflow implements cloudflow.CloudflowServer.
func (svc *templateService) DescribeWorkflow(ctx context.Context, workflow_uri string, tx rdb.Tx) (po.StateMachine, error) {
	var workflow po.StateMachine
	var err error
	tx, maker := svc.metadb.NewTxMaker(tx)
	defer maker.Close(&err)

	err = tx.Where(po.Activity{URI: workflow_uri}).Take(&workflow).Error
	return workflow, err
}

// ListActivities implements cloudflow.CloudflowServer.
func (svc *templateService) ListActivities(ctx context.Context, req vo.ListActivitiesRequest) (vo.ListActivitiesResponse, error) {

	var activitys []po.Activity
	var err error
	var count int64
	var resp vo.ListActivitiesResponse

	limit, offset := req.PageRequest.Limit()
	tx, maker := svc.metadb.NewTxMaker(nil)
	defer maker.Close(&err)
	tx = tx.Model(new(po.Activity))
	err = tx.Count(&count).Error
	if err != nil {
		return resp, err
	}

	err = tx.Offset(limit).Offset(offset).Order("name asc").Find(&activitys).Error
	if err != nil {
		return resp, err
	}

	resp.Activities = activitys
	resp.PageResponse = req.PageRequest.Response(count)
	return resp, err

}

// ListWorkflows implements cloudflow.CloudflowServer.
func (svc *templateService) ListWorkflows(ctx context.Context, req vo.ListWofkflowsRequest) (vo.ListWorkflowsResponse, error) {

	var workerflows []po.StateMachine
	var err error
	var count int64
	var resp vo.ListWorkflowsResponse

	limit, offset := req.PageRequest.Limit()
	tx, maker := svc.metadb.NewTxMaker(nil)
	defer maker.Close(&err)

	tx = tx.Model(new(po.StateMachine))
	err = tx.Count(&count).Error
	if err != nil {
		return resp, err
	}

	err = tx.Offset(limit).Offset(offset).Order("name asc").Find(&workerflows).Error
	if err != nil {
		return resp, err
	}

	resp.Workflows = workerflows
	resp.PageResponse = req.PageRequest.Response(count)
	return resp, err
}
