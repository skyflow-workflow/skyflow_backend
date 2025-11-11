package template

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"testing"

	"github.com/dolthub/go-mysql-server/server"
	"github.com/skyflow-workflow/skyflow_backbend/mock"
)

var (
	mysqlSource         = "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8&parseTime=True&loc=Local"
	testDB              *sql.DB
	testDBServer        *server.Server
	testMysqlServerOnce sync.Once
	myTemplateService   *templateService
)

func TestMain(m *testing.M) {

	var err error

	mockDBClient := mock.GetMockDBClient()
	myTemplateService := NewTemplateService(mockDBClient)
	err = myTemplateService.SyncSchema(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	// 运行测试
	code := m.Run()
	os.Exit(code)
}
