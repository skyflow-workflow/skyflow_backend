package lock

import (
	"fmt"
	"testing"
	"time"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/mock"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
)

var repo *rdb.DBClient

func init() {
	var err error

	dbclient := mock.GetMockDBClient()
	err = dbclient.SyncTables([]interface{}{new(po.ExecutionShade), new(po.Step)})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func AcquireLock(id int, prefix string) {

	var err error
	lockservice := NewDBLockService(repo)
	lock := lockservice.LockExecution(id)
	err = lock.Lock()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 持有锁一段时间
	defer lock.Unlock()

	fmt.Printf("--> %s sleep 5s in lock \n", prefix)
	time.Sleep(1 * time.Second)
	fmt.Printf("--> %s sleep down \n", prefix)
}

func TestLock(t *testing.T) {

	exeid := 5

	for i := 0; i < 4; i++ {
		go AcquireLock(exeid, fmt.Sprintf("lockid-%d", i))
	}
	time.Sleep(10 * time.Second)
}
