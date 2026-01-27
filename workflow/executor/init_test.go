package executor

import (
	"os"
	"sync"
	"testing"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backend/mock"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
)

var (
	testDBClient       *rdb.DBClient
	testQueue          = queue.MockInnerQueue{}
	myExecutionService ExecutionService
	testTestEnvInit    sync.Once
)

func TestMain(m *testing.M) {

	testTestEnvInit.Do(setupTestEnv)

	// 运行测试
	code := m.Run()
	os.Exit(code)
}

func setupTestEnv() {
	testDBClient = mock.GetMockDBClient()
	// 删除所有表
	// testDBClient.DropTables(po.GetExecutionTables())
	testDBClient.SyncTables(po.GetExecutionTables())
	_ = mock.GetMockKafkaMQ()
	myExecutionService = NewExecutionService(testDBClient, testQueue, nil)
}
