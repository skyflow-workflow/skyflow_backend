package template

import (
	"database/sql"
	"os"
	"sync"
	"testing"

	"github.com/dolthub/go-mysql-server/server"
	"github.com/mmtbak/microlibrary/rdb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	mysqlSource         = "root:root@tcp(127.0.0.1:3306)/testdb?charset=utf8&parseTime=True&loc=Local"
	testDB              *sql.DB
	testDBServer        *server.Server
	testMysqlServerOnce sync.Once
	testDBClient        *rdb.DBClient
)

func TestMain(m *testing.M) {

	testMysqlServerOnce.Do(setupTestDB)
	defer func() {
		testDBServer.Close()
		testDB.Close()
	}()
	// 运行测试
	code := m.Run()
	os.Exit(code)
}

func getTestDB() *sql.DB {
	return testDB
}

func getTestDBClient() *rdb.DBClient {
	return testDBClient
}

// setupTestDB 初始化测试用的内存MySQL数据库
func setupTestDB() {

	var err error
	testDB, err = sql.Open("mysql", mysqlSource)
	if err != nil {
		panic(err)
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: getTestDB(),
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	testDBClient = (&rdb.DBClient{}).WithDB(gormDB)
}
