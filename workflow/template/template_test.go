package template

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateNamespace(t *testing.T) {

	var err error
	myTemplateService := NewTemplateService(getTestClient())
	err = myTemplateService.SyncSchema(context.Background(), nil)
	assert.Equal(t, err, nil)
	ns, err := myTemplateService.CreateNamespace(context.Background(), vo.CreateNamespaceRequest{
		Name:    "unittest",
		Comment: "unittest",
	}, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, ns.Data.Name, "unittest")
	fmt.Println(ns)
}

func TestCreateWorkflow(t *testing.T) {

	var err error
	ctx := context.Background()

	myTemplateService := NewTemplateService(getTestClient())
	ns, err := myTemplateService.CreateNamespace(ctx, vo.CreateNamespaceRequest{
		Name:    "unittest",
		Comment: "unittest",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, ns.Data.Name, "unittest")
	assert.Equal(t, ns.Data.Comment, "unittest")

	activity, err := myTemplateService.CreateActivity(ctx, vo.CreateActivityRequest{
		ActivityName: "add",
		Namespace:    "unittest",
		Comment:      "add for add ",
		Parameters:   `{"a": "int", "b": "int"}`,
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	expectedActivity := po.Activity{
		Name:        "add",
		URI:         "activity:unittest/add",
		Comment:     "add for add ",
		Parameters:  `{"a": "int", "b": "int"}`,
		Status:      "Enable",
		NamespaceID: ns.Data.ID,
	}
	assert.Equal(t, activity.Data.Parameters, expectedActivity.Parameters)
	assert.Equal(t, activity.Data.URI, expectedActivity.URI)
	assert.Equal(t, activity.Data.Name, expectedActivity.Name)
	assert.Equal(t, activity.Data.Comment, expectedActivity.Comment)
	assert.Equal(t, activity.Data.Status, expectedActivity.Status)
	assert.Equal(t, activity.Data.NamespaceID, expectedActivity.NamespaceID)

	// 如果参数为空，则设置为空对象
	workflow, err := myTemplateService.CreateStateMachine(ctx, vo.CreateStateMachineRequest{
		StateMachineName: "pass_task",
		Namespace:        "unittest",
		Comment:          "add for add",
		// Definition:   definition,
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, workflow.Data.URI, "statemachine:unittest/pass_task")
	assert.Equal(t, workflow.Data.Name, "pass_task")
	assert.Equal(t, workflow.Data.Comment, "add for add")
	assert.Equal(t, workflow.Data.Status, "Enable")
	assert.Equal(t, workflow.Data.NamespaceID, ns.Data.ID)

	activityquery, err := myTemplateService.DescribeActivity(ctx, expectedActivity.URI, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, activityquery.Parameters, `{"a": "int", "b": "int"}`)

	assert.Equal(t, activityquery.Name, "add")
	assert.Equal(t, activityquery.Comment, "add for add")

}

func TestCreateOrUpdateWorkflow(t *testing.T) {

	var err error
	ctx := context.Background()

	myTemplateService := NewTemplateService(getTestClient())
	// 创建
	ns, err := myTemplateService.CreateOrUpdateNamespace(ctx, vo.CreateNamespaceRequest{
		Name:    "unittest",
		Comment: "unittest",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, ns.Data.Name, "unittest")
	assert.Equal(t, ns.Data.Comment, "unittest")

	// 更新
	ns, err = myTemplateService.CreateOrUpdateNamespace(ctx, vo.CreateNamespaceRequest{
		Name:    "unittest",
		Comment: "unittest update",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, ns.Data.Name, "unittest")
	assert.Equal(t, ns.Data.Comment, "unittest update")

	activity, err := myTemplateService.CreateOrUpdateActivity(ctx, vo.CreateActivityRequest{
		ActivityName: "add",
		Namespace:    "unittest",
		Comment:      "usage for activity add",
		Parameters:   `{"a": "int", "b": "int"}`,
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	expectedActivity := po.Activity{
		Name:        "add",
		URI:         "activity:unittest/add",
		Comment:     "usage for activity add",
		Parameters:  `{"a": "int", "b": "int"}`,
		Status:      "Enable",
		NamespaceID: ns.Data.ID,
	}
	assert.Equal(t, activity.Data.Parameters, expectedActivity.Parameters)
	assert.Equal(t, activity.Data.URI, expectedActivity.URI)
	assert.Equal(t, activity.Data.Name, expectedActivity.Name)
	assert.Equal(t, activity.Data.Comment, expectedActivity.Comment)
	assert.Equal(t, activity.Data.Status, expectedActivity.Status)
	assert.Equal(t, activity.Data.NamespaceID, expectedActivity.NamespaceID)

	activity, err = myTemplateService.CreateOrUpdateActivity(ctx, vo.CreateActivityRequest{
		ActivityName: "add",
		Namespace:    "unittest",
		Comment:      "usage for activity add update",
		Parameters:   `{"a": "int", "b": "int"}`,
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	expectedActivity = po.Activity{
		Name:        "add",
		URI:         "activity:unittest/add",
		Comment:     "usage for activity add update",
		Parameters:  `{"a": "int", "b": "int"}`,
		Status:      "Enable",
		NamespaceID: ns.Data.ID,
	}
	assert.Equal(t, activity.Data.Parameters, expectedActivity.Parameters)
	assert.Equal(t, activity.Data.URI, expectedActivity.URI)
	assert.Equal(t, activity.Data.Name, expectedActivity.Name)
	assert.Equal(t, activity.Data.Comment, expectedActivity.Comment)
	assert.Equal(t, activity.Data.Status, expectedActivity.Status)
	assert.Equal(t, activity.Data.NamespaceID, expectedActivity.NamespaceID)

	// 创建
	workflow, err := myTemplateService.CreateOrUpdateStateMachine(ctx, vo.CreateStateMachineRequest{
		StateMachineName: "pass_task",
		Namespace:        "unittest",
		Comment:          "usage for workflow pass_task",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, workflow.Data.URI, "statemachine:unittest/pass_task")
	assert.Equal(t, workflow.Data.Name, "pass_task")
	assert.Equal(t, workflow.Data.Comment, "usage for workflow pass_task")
	assert.Equal(t, workflow.Data.Status, "Enable")
	assert.Equal(t, workflow.Data.NamespaceID, ns.Data.ID)
	// 更新
	workflow, err = myTemplateService.CreateOrUpdateStateMachine(ctx, vo.CreateStateMachineRequest{
		StateMachineName: "pass_task",
		Namespace:        "unittest",
		Comment:          "usage for workflow pass_task update",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, workflow.Data.URI, "statemachine:unittest/pass_task")
	assert.Equal(t, workflow.Data.Name, "pass_task")
	assert.Equal(t, workflow.Data.Comment, "usage for workflow pass_task update")
	assert.Equal(t, workflow.Data.Status, "Enable")
	assert.Equal(t, workflow.Data.NamespaceID, ns.Data.ID)

	activityquery, err := myTemplateService.DescribeActivity(ctx, expectedActivity.URI, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, activityquery.Parameters, `{"a": "int", "b": "int"}`)

	assert.Equal(t, activityquery.Name, "add")
	assert.Equal(t, activityquery.Comment, "usage for activity add update")

}
