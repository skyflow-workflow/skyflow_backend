package lock

import (
	"fmt"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
)

type DBLockService struct {
	client *rdb.DBClient
}

type DBLock struct {
	service  *DBLockService
	id       int
	tx       rdb.Tx
	locktype string
}

func NewDBLockService(repo *rdb.DBClient) *DBLockService {
	ls := &DBLockService{
		client: repo,
	}
	return ls
}

func (db *DBLockService) LockExecution(id int) Lock {

	lock := &DBLock{
		service:  db,
		id:       id,
		locktype: LockTypes.Execution,
	}
	return lock
}
func (db *DBLockService) LockStep(id int) Lock {

	lock := &DBLock{
		service:  db,
		id:       id,
		locktype: LockTypes.Step,
	}
	return lock
}

func (lock *DBLock) Lock() error {

	var err error
	tx := lock.service.client.NewTx()
	var lockExe po.ExecutionShade
	var lockstep po.Step
	tx.Begin()
	switch lock.locktype {
	case LockTypes.Execution:
		err = rdb.ForUpdate(tx).Select("id").Take(&lockExe, lock.id).Error
	case LockTypes.Step:
		err = rdb.ForUpdate(tx).Select("id").Take(&lockstep, lock.id).Error
	default:
		err = fmt.Errorf("unsupported lock resource type  '%s'", lock.locktype)
	}
	if err != nil {
		return err
	}
	lock.tx = tx
	return nil

}

// UnLock unlock state
func (lock *DBLock) Unlock() {

	if lock.tx == nil {
		return
	}
	lock.tx.Rollback()

}

func (lock *DBLock) GetTx() rdb.Tx {
	return lock.tx
}
