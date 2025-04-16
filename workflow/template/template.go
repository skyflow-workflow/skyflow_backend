// package template
package template

import (
	"context"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/gen/pb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/pberror"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// TemplateService template management service
type TemplateService = *templateService

type templateService struct {
	dbclient *rdb.DBClient
}

func NewTemplateService(dbclient *rdb.DBClient) TemplateService {
	svc := &templateService{
		dbclient: dbclient,
	}
	return svc
}

func (svc *templateService) SyncSchema(ctx context.Context, tx rdb.Tx) error {

	var err error

	tx, maker := svc.dbclient.NewTxMaker(tx)
	defer maker.Close(&err)

	err = tx.AutoMigrate(po.GetTemplateTables()...)
	if err != nil {
		return err
	}

	return nil

}

func (svc *templateService) CreateNamespace(ctx context.Context, req vo.CreateNamespaceRequest, tx rdb.Tx) (vo.CreateNamespaceResponse, error) {

	var ns po.Namespace
	var err error

	tx, maker := svc.dbclient.NewTxMaker(tx)
	defer maker.Close(&err)

	err = tx.Where("name = ?", req.Name).First(&ns).Error
	if err != nil && !rdb.IsErrRecordNotFound(err) {
		return vo.CreateNamespaceResponse{}, err
	}

	if ns.ID != 0 {
		// namespace already exists
		err = pberror.NewPBError(int32(pb.ErrorCode_NameSpaceAlreadyExists), "namespace: "+req.Name+" already exists")
		return vo.CreateNamespaceResponse{}, err
	}

	new_ns := po.Namespace{
		Name:    req.Name,
		Comment: req.Comment,
	}

	err = tx.Create(&new_ns).Error
	if err != nil {
		tx.Commit()
	}
	return vo.CreateNamespaceResponse{
		Data: new_ns,
	}, err
}

func (svc *templateService) DescribeNamespace(ctx context.Context, name string, tx rdb.Tx) (po.Namespace, error) {
	var ns po.Namespace
	var err error

	tx, maker := svc.dbclient.NewTxMaker(tx)
	defer maker.Close(&err)

	err = tx.Where("name = ?", name).First(&ns).Error
	if err != nil && rdb.IsErrRecordNotFound(err) {
		// namespace not found
		err = pberror.NewPBError(int32(pb.ErrorCode_NameSpaceDoesNotExist), "namespace: "+name+" not found")
		return ns, err
	}
	if err != nil {
		return ns, err
	}
	return ns, nil
}

func (svc *templateService) DeleteNamespace(ctx context.Context, req vo.DeleteNamespaceRequest, tx rdb.Tx) error {
	var ns po.Namespace
	var err error

	tx, maker := svc.dbclient.NewTxMaker(tx)
	defer maker.Close(&err)

	// 先查询是否存在
	err = tx.Where("name = ?", req.Name).First(&ns).Error
	if err != nil && rdb.IsErrRecordNotFound(err) {
		// namespace not found
		err = pberror.NewPBError(int32(pb.ErrorCode_NameSpaceDoesNotExist), "namespace: "+req.Name+" not found")
		return err
	}
	if err != nil {
		// other error
		return err
	}

	// 删除
	err = tx.Where(po.Namespace{ID: ns.ID}).Delete(&po.Namespace{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (svc *templateService) CreateOrUpdateNamespace(ctx context.Context, req vo.CreateNamespaceRequest, tx rdb.Tx) (vo.CreateNamespaceResponse, error) {
	var ns po.Namespace
	var err error

	tx, maker := svc.dbclient.NewTxMaker(tx)
	defer maker.Close(&err)

	// 先查询是否存在
	err = tx.Where("name = ?", req.Name).First(&ns).Error
	if err != nil && !rdb.IsErrRecordNotFound(err) {
		return vo.CreateNamespaceResponse{}, err
	}

	newns := po.Namespace{
		Name:    req.Name,
		Comment: req.Comment,
	}
	// 如果存在则更新，否则创建,
	if ns.ID != 0 {
		// 为什么不用Save， Save会更新所有字段，而我们只需要更新部分字段
		err = tx.Where(po.Namespace{ID: ns.ID}).Updates(&newns).Error
		return vo.CreateNamespaceResponse{
			Data: newns,
		}, err
	}
	err = tx.Create(&newns).Error
	if err != nil && rdb.IsErrRecordNotFound(err) {
		return vo.CreateNamespaceResponse{}, err
	}

	return vo.CreateNamespaceResponse{
		Data: newns,
	}, err
}

// ListNamespaces implements skyflow.SkyflowServer.
func (svc *templateService) ListNamespaces(ctx context.Context, req vo.ListNamespacesRequest) (vo.ListNamespacesResponse, error) {

	var err error
	var count int64
	var nss = []po.Namespace{}
	var resp vo.ListNamespacesResponse
	tx, maker := svc.dbclient.NewTxMaker(nil)
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

// CreateActivity implements skyflow.SkyflowServer.
func (svc *templateService) CreateActivity(ctx context.Context, req vo.CreateActivityRequest, tx rdb.Tx) (vo.CreateActivityResponse, error) {

	var activity po.Activity
	var err error

	activity_uri := parser.GenerateActivityURI(req.Namespace, req.ActivityName)

	tx, maker := svc.dbclient.NewTxMaker(tx)
	defer maker.Close(&err)

	queryActivity := po.Activity{
		URI: activity_uri,
	}

	err = tx.Model(new(po.Activity)).Where(queryActivity).Select("id").First(&activity).Error
	if err != nil && !rdb.IsErrRecordNotFound(err) {
		return vo.CreateActivityResponse{}, err
	}

	namespace, err := svc.DescribeNamespace(ctx, req.Namespace, tx)
	if err != nil {
		return vo.CreateActivityResponse{}, err
	}

	newActivity := po.Activity{
		NamespaceID: namespace.ID,
		Name:        req.ActivityName,
		Comment:     req.Comment,
		URI:         activity_uri,
		Status:      ActivityStatus.Enable,
	}

	if activity.ID != 0 {
		err = tx.Where(po.Activity{ID: activity.ID}).Updates(&newActivity).Error
	} else {
		err = tx.Create(&newActivity).Error
	}

	if err != nil {
		return vo.CreateActivityResponse{}, err
	}

	return vo.CreateActivityResponse{
		Data: newActivity,
	}, nil
}

// ListActivities implements skyflow.SkyflowServer.
func (svc *templateService) ListActivities(ctx context.Context, req vo.ListActivitiesRequest) (vo.ListActivitiesResponse, error) {

	var activitys []po.Activity
	var err error
	var count int64
	var resp vo.ListActivitiesResponse

	limit, offset := req.PageRequest.Limit()
	tx, maker := svc.dbclient.NewTxMaker(nil)
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

func (svc *templateService) CreateOrUpdateActivity(ctx context.Context, req vo.CreateActivityRequest, tx rdb.Tx) (vo.CreateActivityResponse, error) {

	var activity po.Activity
	var err error

	activity_uri := parser.GenerateActivityURI(req.Namespace, req.ActivityName)

	tx, maker := svc.dbclient.NewTxMaker(tx)
	defer maker.Close(&err)

	queryActivity := po.Activity{
		URI: activity_uri,
	}

	err = tx.Model(new(po.Activity)).Where(queryActivity).Select("id").First(&activity).Error
	if err != nil && !rdb.IsErrRecordNotFound(err) {
		return vo.CreateActivityResponse{}, err
	}

	namespace, err := svc.DescribeNamespace(ctx, req.Namespace, tx)
	if err != nil {
		return vo.CreateActivityResponse{}, err
	}

	newActivity := po.Activity{
		NamespaceID: namespace.ID,
		Name:        req.ActivityName,
		Comment:     req.Comment,
		URI:         activity_uri,
		Status:      ActivityStatus.Enable,
	}

	if activity.ID != 0 {
		err = tx.Where(po.Activity{ID: activity.ID}).Updates(&newActivity).Error
	} else {
		err = tx.Create(&newActivity).Error
	}

	if err != nil {
		return vo.CreateActivityResponse{}, err
	}
	return vo.CreateActivityResponse{
		Data: newActivity,
	}, nil
}

// DeleteActivity implements skyflow.SkyflowServer.
func (svc *templateService) DeleteActivity(ctx context.Context, req vo.DeleteActivityRequest, tx rdb.Tx) error {
	var activity po.Activity
	var err error

	tx, maker := svc.dbclient.NewTxMaker(tx)
	defer maker.Close(&err)

	// 先查询是否存在
	err = tx.Where("name = ?", req.ActivityURI).First(&activity).Error
	if err != nil && rdb.IsErrRecordNotFound(err) {
		// activity not found
		err = pberror.NewPBError(int32(pb.ErrorCode_ActivityDoesNotExist), "activity: "+req.ActivityURI+" not found")
		return err
	}
	if err != nil {
		return err
	}

	// 删除
	err = tx.Where(po.Activity{ID: activity.ID}).Delete(&po.Activity{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// CreateStateMachine implements skyflow.SkyflowServer.
func (svc *templateService) CreateStateMachine(ctx context.Context, req vo.CreateStateMachineRequest, tx rdb.Tx) (vo.CreateStateMachineResponse, error) {

	var err error
	var sm po.StateMachine
	flow, err := parser.ParseStateMachine(req.Definition)
	if err != nil {
		return vo.CreateStateMachineResponse{}, err
	}

	workflowUri := parser.GenerateStateMachineURI(req.Namespace, req.StateMachineName)

	tx, maker := svc.dbclient.NewTxMaker(tx)
	defer maker.Close(&err)

	queryWorkflow := po.StateMachine{
		URI: workflowUri,
	}

	err = tx.Model(new(po.StateMachine)).Where(queryWorkflow).Select("id").First(&sm).Error
	if err != nil && !rdb.IsErrRecordNotFound(err) {
		return vo.CreateStateMachineResponse{}, err
	}

	namespace, err := svc.DescribeNamespace(ctx, req.Namespace, tx)
	if err != nil {
		return vo.CreateStateMachineResponse{}, err
	}

	newsm := po.StateMachine{
		NamespaceID: namespace.ID,
		Name:        req.StateMachineName,
		Type:        flow.Type,
		Comment:     req.Comment,
		URI:         workflowUri,
		Definition:  req.Definition,
		Status:      ActivityStatus.Enable,
	}
	if sm.ID != 0 {
		err = tx.Where(po.StateMachine{ID: sm.ID}).Updates(&newsm).Error
		return vo.CreateStateMachineResponse{
			Data: newsm,
		}, err
	}

	err = tx.Create(&newsm).Error

	if err != nil && rdb.IsErrRecordNotFound(err) {
		return vo.CreateStateMachineResponse{}, err
	}
	return vo.CreateStateMachineResponse{
		Data: newsm,
	}, nil
}

// DescribeActivity implements skyflow.SkyflowServer.
func (svc *templateService) DescribeActivity(ctx context.Context, activityUri string, tx rdb.Tx) (po.Activity, error) {
	var activity po.Activity
	var err error
	tx, maker := svc.dbclient.NewTxMaker(tx)
	defer maker.Close(&err)

	err = tx.Where(po.Activity{URI: activityUri}).Take(&activity).Error
	return activity, err
}

// DescribeWorkflow implements skyflow.SkyflowServer.
func (svc *templateService) DescribeWorkflow(ctx context.Context, stateMachineUri string, tx rdb.Tx) (po.StateMachine, error) {
	var workflow po.StateMachine
	var err error
	tx, maker := svc.dbclient.NewTxMaker(tx)
	defer maker.Close(&err)

	err = tx.Where(po.Activity{URI: stateMachineUri}).Take(&workflow).Error
	return workflow, err
}

// ListWorkflows implements skyflow.SkyflowServer.
func (svc *templateService) ListStateMachines(ctx context.Context, req vo.ListStateMachinesRequest, tx rdb.Tx) (vo.ListStateMachinesResponse, error) {

	var sms []po.StateMachine
	var err error
	var count int64
	var resp vo.ListStateMachinesResponse

	limit, offset := req.PageRequest.Limit()
	tx, maker := svc.dbclient.NewTxMaker(tx)
	defer maker.Close(&err)

	tx = tx.Model(new(po.StateMachine))
	err = tx.Count(&count).Error
	if err != nil {
		return resp, err
	}

	err = tx.Offset(limit).Offset(offset).Order("name asc").Find(&sms).Error
	if err != nil {
		return resp, err
	}

	resp.StateMachines = sms
	resp.PageResponse = req.PageRequest.Response(count)
	return resp, err
}
