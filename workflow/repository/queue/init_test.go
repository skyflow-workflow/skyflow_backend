package queue

import (
	"database/sql"
	"os"
	"sync"
	"testing"

	"github.com/dolthub/go-mysql-server/server"
	"github.com/mmtbak/microlibrary/rdb"
)

var (
	testDB              *sql.DB
	testDBServer        *server.Server
	testMysqlServerOnce sync.Once
	testDBClient        *rdb.DBClient
)

func TestMain(m *testing.M) {

	// 运行测试
	code := m.Run()
	os.Exit(code)
}
