package executor

import (
	"os"
	"sync"
	"testing"

	"github.com/mmtbak/microlibrary/mq"
	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/mock"
)

var (
	testDBClient       *rdb.DBClient
	testKafkaMQ        *mq.KafkaMessageQueue
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
	_ = mock.GetMockKafkaMQ()
	myExecutionService = NewExecutionService(testDBClient, nil, nil)
}
