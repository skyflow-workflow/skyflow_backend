package executor

//query data for executor

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backend/workflow/vo"
)

// QueryStepByID
func (svc *executionService) QueryStepByID(step_id int, fields []string, session rdb.Tx) (*po.Step, error) {

	var dbStep po.Step
	var err error
	slog.Debug("QueryStepByID",
		"step_id", step_id,
		"fields", fields,
	)
	// 增加控制session
	tx, maker := svc.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)
	// 判断字段
	if len(fields) != 0 {
		tx = tx.Select(fields)
	}
	err = tx.Take(&dbStep, step_id).Error
	if rdb.IsErrRecordNotFound(err) {
		return nil, fmt.Errorf("%w: %d", vo.ErrorStepNotFound, step_id)
	}
	return &dbStep, err
}

// QueryStepIDByTaskToken 使用tasktoken 查询step id,
func (svc *executionService) QueryStepIDByTaskToken(tasktoken string, session rdb.Tx) (step_id int, err error) {

	// 增加控制session
	tx, maker := svc.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)

	err = tx.Model(new(po.TaskToken)).Where(po.TaskToken{Token: tasktoken, IsDeleted: false}).Select("step_id").Take(&step_id).Error
	if err != nil {
		if rdb.IsErrRecordNotFound(err) {
			err = fmt.Errorf("%w: %s", vo.ErrorTaskTokenNotFound, tasktoken)
			return
		}
		return
	}
	return step_id, err
}

// QueryExecutionByID
func (svc *executionService) QueryExecutionByID(execution_id int, fields []string, session rdb.Tx) (*po.Execution, error) {
	var dbExe po.Execution
	var err error
	// 增加控制session
	tx, maker := svc.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)
	// 判断字段
	if len(fields) != 0 {
		tx = tx.Select(fields)
	}
	err = tx.Take(&dbExe, execution_id).Error
	if rdb.IsErrRecordNotFound(err) {
		return nil, fmt.Errorf("%w: execution id %d", vo.ErrorExecutionNotFound, execution_id)
	}
	return &dbExe, err
}

// QueryExecutionByUUID
func (svc *executionService) QueryExecutionByUUID(uuid string, fields []string, session rdb.Tx) (*po.Execution, error) {
	var dbExe po.Execution
	var err error
	// 增加控制session
	tx, maker := svc.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)
	// 判断字段
	if len(fields) != 0 {
		tx = tx.Select(fields)
	}
	err = tx.Where(po.Execution{UUID: uuid}).Take(&dbExe).Error
	if rdb.IsErrRecordNotFound(err) {
		return nil, fmt.Errorf("%w: execution uuid '%s' ", vo.ErrorExecutionNotFound, uuid)
	}
	return &dbExe, err
}

// QueryStepGroupByStepID
func (svc *executionService) QueryStepGroupByStepID(step_id int, session rdb.Tx) (*po.StepGroup, error) {
	var dbStepgroup po.StepGroup
	var err error

	// 增加控制session
	tx, maker := svc.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)
	// 判断字段
	// stategroup 字段少，不需要判断field
	err = tx.Where(po.StepGroup{StepID: step_id}).Take(&dbStepgroup).Error
	if rdb.IsErrRecordNotFound(err) {
		return nil, fmt.Errorf("%w: %d", vo.ErrorStepGroupNotFound, step_id)
	}
	return &dbStepgroup, err
}

// QueryStepGroupByStepID
func (svc *executionService) QueryStepGroupBySubGroupID(execution_id int, subgroup_id int, session rdb.Tx) (*po.StepGroup, error) {
	var dbStepgroup po.StepGroup
	var err error

	// 增加控制session
	tx, maker := svc.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)
	// 判断字段

	cond := po.StepGroup{
		ExecutionID: execution_id,
		SubGroupID:  subgroup_id,
	}

	// stategroup 字段少，不需要判断field
	err = tx.Where(cond).Take(&dbStepgroup).Error
	return &dbStepgroup, err
}

func (svc *executionService) QueryStepUserData(step_id int, fields []string, session rdb.Tx) (*po.UserStepData, error) {
	var dbUserData po.UserStepData
	var err error
	// 增加控制session
	tx, maker := svc.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)
	// 判断字段
	if len(fields) != 0 {
		tx = tx.Select(fields)
	}
	err = tx.Take(&dbUserData, step_id).Error
	if rdb.IsErrRecordNotFound(err) {
		return nil, fmt.Errorf("%w: %d", vo.ErrorUserStepDataNotFound, step_id)
	}
	return &dbUserData, err
}

// ListExecutions 获得execution 列表
func (svc *executionService) ListExecutions(req vo.ListExecutionsRequest) (vo.ListExecutionsResponse, error) {

	var err error
	var dbExecutions = []po.Execution{}
	var count int64
	var resp vo.ListExecutionsResponse
	// 新建事务
	tx, maker := svc.GetMetaDB().NewTxMaker(nil)
	defer maker.Close(&err)

	if req.Status != "" {
		tx = tx.Where("status = ?", req.Status)
	}
	if req.WorkflowURI != "" {
		tx = tx.Where("uri = ?", req.WorkflowURI)
	}
	title := strings.TrimSpace(req.Title)
	if title != "" {
		tx = tx.Where("title like ? ", "%"+title+"%")
	}
	if len(req.ExecutionUUIDs) != 0 {
		tx = tx.Where("uuid in ?", req.ExecutionUUIDs)
	}

	limit, offset := req.PageRequest.Limit()

	tx = tx.Model(new(po.Execution))
	// 查询总数
	err = tx.Count(&count).Error
	if err != nil {
		return resp, err
	}
	// 查询数据
	fields := append(
		ExecutionFields.L2,
		ExecutionFieldNames.Title,
		ExecutionFieldNames.URI,
		ExecutionFieldNames.CreateTime,
		ExecutionFieldNames.StartTime,
		ExecutionFieldNames.FinishTime,
		ExecutionFieldNames.UpdateTime,
	)
	err = tx.Limit(limit).Offset(offset).Select(fields).Order("id DESC").Find(&dbExecutions).Error
	if err != nil {
		return resp, err
	}
	resp = vo.ListExecutionsResponse{
		Executions:   dbExecutions,
		PageResponse: req.PageRequest.Response(count),
	}
	return resp, nil
}

// DescribeExecutionBone 获得一个Execution 的bone结构
func (svc *executionService) DescribeExecutionBone(ctx context.Context, req vo.DescribeExecutionBoneRequest) (
	resp vo.DescribeExecutionBoneResponse, err error) {

	var dbExecution *po.Execution
	var dbSteps = []po.Step{}
	var dbStepgroups = []po.StepGroup{}
	var bone interface{}

	tx, maker := svc.GetMetaDB().NewTxMaker(nil)
	defer maker.Close(&err)

	// 获得Execution
	dbExecution, err = svc.QueryExecutionByID(req.ExecutionID, append(ExecutionFields.L1, "header"), tx)
	if err != nil {
		return
	}

	// 获得Steps列表
	err = tx.Where(po.Step{ExecutionID: req.ExecutionID}).
		Select(StepFields.L2, StepFieldNames.Resource).
		Find(&dbSteps).Error
	if err != nil {
		return
	}
	// 获得StepGroups列表
	err = tx.Where(po.StepGroup{ExecutionID: req.ExecutionID}).Find(&dbStepgroups).Error
	if err != nil {
		return
	}
	// 已经查询完了， 可以结束数据了。后面就是计算了
	tx.Commit()

	exeObj, err := NewExecutionFromData(dbExecution, svc)
	if err != nil {
		return
	}
	err = exeObj.FullInit()
	if err != nil {
		return
	}
	bone, err = exeObj.GetBone()
	if err != nil {
		return
	}
	resp = vo.DescribeExecutionBoneResponse{
		Bone: bone,
	}
	return resp, nil

}
