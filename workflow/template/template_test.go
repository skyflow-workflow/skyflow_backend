package template

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
	"github.com/dolthub/vitess/go/vt/proto/query"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbName    = "mydb"
	tableName = "mytable"
	address   = "localhost"
	port      = 3306
)

func TestWithGoMySQLServer(t *testing.T) {
	// 1. 创建内存数据库
	db := memory.NewDatabase("testdb")

	pro := memory.NewDBProvider(db)
	session := memory.NewSession(sql.NewBaseSession(), pro)
	ctx := sql.NewContext(context.Background(), sql.WithSession(session))
	// 2. 创建表结构
	table := memory.NewTable(db, tableName, sql.NewPrimaryKeySchema(sql.Schema{
		{Name: "name", Type: types.Text, Nullable: false, Source: tableName, PrimaryKey: true},
		{Name: "email", Type: types.Text, Nullable: false, Source: tableName, PrimaryKey: true},
		{Name: "phone_numbers", Type: types.JSON, Nullable: false, Source: tableName},
		{Name: "created_at", Type: types.MustCreateDatetimeType(query.Type_DATETIME, 6), Nullable: false, Source: tableName},
	}), db.GetForeignKeyCollection())
	db.AddTable(tableName, table)

	creationTime := time.Unix(0, 1667304000000001000).UTC()
	_ = table.Insert(ctx, sql.NewRow("Jane Deo", "janedeo@gmail.com", types.MustJSON(`["556-565-566", "777-777-777"]`), creationTime))
	_ = table.Insert(ctx, sql.NewRow("Jane Doe", "jane@doe.com", types.MustJSON(`[]`), creationTime))
	_ = table.Insert(ctx, sql.NewRow("John Doe", "john@doe.com", types.MustJSON(`["555-555-555"]`), creationTime))
	_ = table.Insert(ctx, sql.NewRow("John Doe", "johnalt@doe.com", types.MustJSON(`[]`), creationTime))

	// 3. 创建引擎
	engine := sqle.NewDefault(memory.NewDBProvider(db))
	ctx = sql.NewContext(context.Background(), sql.WithSession(session))
	ctx.SetCurrentDatabase(dbName)

	// 4. 配置服务器
	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", address, port),
	}

	// 5. 启动服务器
	srv, err := server.NewServer(config, engine, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		if err := srv.Start(); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()
	defer srv.Close()

	// 6. 等待服务器启动
	time.Sleep(500 * time.Millisecond)
	select {}

	// // 7. 连接测试
	// dsn := "root:@tcp(localhost:3306)/testdb"
	// dbConn, err := sql.Open("mysql", dsn)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer dbConn.Close()

	// // 8. 执行测试查询
	// _, err = dbConn.Exec("INSERT INTO users (id, name, email) VALUES (1, 'John', 'john@example.com')")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// var name string
	// err = dbConn.QueryRow("SELECT name FROM users WHERE id = 1").Scan(&name)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if name != "John" {
	// 	t.Errorf("Expected John, got %s", name)
	// }
}

var myTemplateService TemplateService

func CreateMySQLServer(t *testing.T) {

}

func TestTiDBMemoryMode(t *testing.T) {

}

func TestCreateNamespace(t *testing.T) {

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
