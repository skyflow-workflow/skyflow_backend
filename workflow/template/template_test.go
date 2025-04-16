package template

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbName  = "testdb"
	address = "127.0.0.1"
	port    = 3306
	dsn     = "root:@tcp(localhost:3306)/testdb?charset=utf8&parseTime=True&loc=Local"
)

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

var myTemplateService TemplateService

func TestCreateNamespace(t *testing.T) {

	mysqlsever, err := CreateMemoryMySQLServer()
	assert.Equal(t, err, nil)
	go func() {
		slog.Info("start server ")
		if err = mysqlsever.Start(); err != nil {
			panic(err)
		}
	}()
	defer mysqlsever.Close()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.Equal(t, err, nil)

	client := (&rdb.DBClient{}).WithDB(db)
	myTemplateService := NewTemplateService(client)
	ns, err := myTemplateService.CreateNamespace(context.Background(), vo.CreateNamespaceRequest{
		Name:    "unittest",
		Comment: "unittest",
	}, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, ns.Data.Name, "unittest")

}

func TestCreateWorkflow(t *testing.T) {

	ctx := context.Background()
	var err error
	db, mock, err := sqlmock.New()
	assert.Equal(t, err, nil)
	// mock statement
	// mock sql "select version()"
	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7.30"))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `mock_tables`").WillReturnResult(sqlmock.NewResult(10, 1))
	mock.ExpectExec("UPDATE `mock_tables`").WithArgs("bob", 20, 10).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	// mock end statement
	gormdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	assert.Equal(t, err, nil)
	client := (&rdb.DBClient{}).WithDB(gormdb)
	myTemplateService = NewTemplateService(client)
	ns, err := myTemplateService.CreateNamespace(ctx, vo.CreateNamespaceRequest{
		Name:    "unittest",
		Comment: "unittest",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, ns.Data.Name, "")
	activity, err := myTemplateService.CreateActivity(ctx, vo.CreateActivityRequest{
		ActivityName: "add",
		Namespace:    "unittest",
		Comment:      "add for add ",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(activity)
	//
	workflow, err := myTemplateService.CreateStateMachine(ctx, vo.CreateStateMachineRequest{
		StateMachineName: "pass_task",
		Namespace:        "unittest",
		Comment:          "add for add ",
		// Definition:   definition,
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(workflow)

	activityquery, err := myTemplateService.DescribeActivity(ctx, activity.Data.URI, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(activityquery)

}
