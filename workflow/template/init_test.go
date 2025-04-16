package template

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbName              = "testdb"
	address             = "127.0.0.1"
	port                = 3306
	dsn                 = fmt.Sprintf("root:@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", address, port, dbName)
	testDB              *sql.DB
	testdbserver        *server.Server
	testMysqlServerOnce sync.Once
	testdbClient        *rdb.DBClient
)

func TestMain(m *testing.M) {

	testMysqlServerOnce.Do(setupTestDB)
	defer func() {
		testdbserver.Close()
		testDB.Close()
	}()
	// 运行测试
	code := m.Run()
	os.Exit(code)
}

func getTestDB() *sql.DB {
	return testDB
}

func getTestClient() *rdb.DBClient {
	return testdbClient
}

func setupTestDB() {

	var err error
	testdbserver, err = CreateMemoryMySQLServer()
	if err != nil {
		panic(err)
	}
	StartMySQLServer(testdbserver)

	testDB, err = sql.Open("mysql", dsn)
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

	testdbClient = (&rdb.DBClient{}).WithDB(gormDB)
	testdbClient.SyncTables(po.GetTemplateTables())
}

func StartMySQLServer(server *server.Server) {
	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(500 * time.Millisecond)
}

func CreateMemoryMySQLServer() (*server.Server, error) {
	db := memory.NewDatabase(dbName)
	db.BaseDatabase.EnablePrimaryKeyIndexes()
	provider := memory.NewDBProvider(db)

	engine := sqle.NewDefault(provider)
	// session := memory.NewSession(sql.NewBaseSession(), pro)
	// ctx := sql.NewContext(context.Background(), sql.WithSession(session))
	// ctx.SetCurrentDatabase(dbName)

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", address, port),
	}
	s, err := server.NewServer(config, engine, memory.NewSessionBuilder(provider), nil)
	return s, err
}
