package executor

//query data for executor

import (
	"fmt"
	"log/slog"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
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
		return nil, fmt.Errorf("%w: execution uuid %d", vo.ErrorExecutionNotFound, execution_id)
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
