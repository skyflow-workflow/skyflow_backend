package lock

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/mock"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
)

var dbClient *rdb.DBClient

func init() {
	var err error

	dbClient = mock.GetMockDBClient()
	err = dbClient.SyncTables([]interface{}{new(po.ExecutionShade), new(po.Step)})
	if err != nil {
		log.Fatal(err)
		return
	}
	err = dbClient.TruncateTables([]interface{}{new(po.ExecutionShade), new(po.Step)})
	if err != nil {
		log.Fatal(err)
		return
	}
}

func AcquireLock(id int, prefix string, wg *sync.WaitGroup) {

	var err error
	defer wg.Done()
	lockservice := NewDBLockService(dbClient)
	lock := lockservice.LockExecution(id)
	err = lock.Lock()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 持有锁一段时间
	defer lock.Unlock()

	log.Printf("--> %s sleep 5s in lock \n", prefix)
	time.Sleep(1 * time.Second)
	log.Printf("--> %s sleep down \n", prefix)
	tx := lock.GetTx()
	log.Printf("--> %s lock is tx : %v \n", prefix, tx != nil)
}

func TestExecutionLock(t *testing.T) {

	var err error
	exeShadeID := 5
	err = dbClient.DB().Create(&po.ExecutionShade{
		ID: exeShadeID,
	}).Error
	assert.Equal(t, err, nil)

	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go AcquireLock(exeShadeID, fmt.Sprintf("lockid-%d", i), &wg)
	}
	wg.Wait()
	// wait for all goroutine to finish
	log.Println("all goroutine finished")
}
