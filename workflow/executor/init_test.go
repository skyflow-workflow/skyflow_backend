package executor

import (
	"os"
	"sync"
	"testing"

	"github.com/mmtbak/microlibrary/mq"
	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/mock"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/config"
)

var (
	testDBClient    *rdb.DBClient
	testKafkaMQ     *mq.KafkaMessageQueue
	myExecutor      = StandardExecutor
	testTestEnvInit sync.Once
)

func TestMain(m *testing.M) {

	testTestEnvInit.Do(setupTestEnv)

	// 运行测试
	code := m.Run()
	os.Exit(code)
}

func setupTestEnv() {
	testDBClient = mock.GetMockDBClient()
	testKafkaMQ = mock.GetMockKafkaMQ()
	myExecutor = NewExecutor(&config.StandardExecutorConfig)
	myExecutor.MetaDB = testDBClient
}
