package template

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var myTemplateService TemplateService

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
	ns, err := myTemplateService.CreateNamespace(ctx, "unittest", "", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, ns.Name, "")
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
	workflow, err := myTemplateService.CreateWorkflow(ctx, vo.CreateWorkflowRequest{
		WorkflowName: "pass_task",
		Namespace:    "unittest",
		Comment:      "add for add ",
		// Definition:   definition,
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(workflow)

	activityquery, err := myTemplateService.DescribeActivity(ctx, activity.URI, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(activityquery)

}
