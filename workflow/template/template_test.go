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
		Name:    "testing_create_namespace",
		Comment: "testing_create_namespace",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, ns.Data.Name, "testing_create_namespace")
	assert.Equal(t, ns.Data.Comment, "testing_create_namespace")

	activity, err := myTemplateService.CreateActivity(ctx, vo.CreateActivityRequest{
		ActivityName: "add",
		Namespace:    "testing_create_namespace",
		Comment:      "testing_create_activity_add_description",
		Parameters:   `{"a": "int", "b": "int"}`,
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	expectedActivity := po.Activity{
		Name:        "add",
		URI:         "activity:testing_create_namespace/add",
		Comment:     "testing_create_activity_add_description",
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
		Namespace:        "testing_create_namespace",
		Comment:          "pass_task description",
		// Definition:   definition,
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, workflow.Data.URI, "statemachine:testing_create_namespace/pass_task")
	assert.Equal(t, workflow.Data.Name, "pass_task")
	assert.Equal(t, workflow.Data.Comment, "pass_task description")
	assert.Equal(t, workflow.Data.Status, "Enable")
	assert.Equal(t, workflow.Data.NamespaceID, ns.Data.ID)

	activityquery, err := myTemplateService.DescribeActivity(ctx, expectedActivity.URI, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, activityquery.Parameters, `{"a": "int", "b": "int"}`)

	assert.Equal(t, activityquery.Name, "add")
	assert.Equal(t, activityquery.Comment, "testing_create_activity_add_description")

}

func TestCreateOrUpdateWorkflow(t *testing.T) {

	var err error
	ctx := context.Background()

	myTemplateService := NewTemplateService(getTestClient())
	// 创建
	ns, err := myTemplateService.CreateOrUpdateNamespace(ctx, vo.CreateNamespaceRequest{
		Name:    "testing_create_or_update_namespace",
		Comment: "testing_create_or_update_namespace",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, ns.Data.Name, "testing_create_or_update_namespace")
	assert.Equal(t, ns.Data.Comment, "testing_create_or_update_namespace")

	// 更新
	ns, err = myTemplateService.CreateOrUpdateNamespace(ctx, vo.CreateNamespaceRequest{
		Name:    "testing_create_or_update_namespace",
		Comment: "testing_create_or_update_namespace update",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, ns.Data.Name, "testing_create_or_update_namespace")
	assert.Equal(t, ns.Data.Comment, "testing_create_or_update_namespace update")

	activity, err := myTemplateService.CreateOrUpdateActivity(ctx, vo.CreateActivityRequest{
		ActivityName: "add",
		Namespace:    "testing_create_or_update_namespace",
		Comment:      "usage for activity add",
		Parameters:   `{"a": "int", "b": "int"}`,
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	expectedActivity := po.Activity{
		Name:        "add",
		URI:         "activity:testing_create_or_update_namespace/add",
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
		Namespace:    "testing_create_or_update_namespace",
		Comment:      "usage for activity add update",
		Parameters:   `{"a": "int", "b": "int"}`,
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	expectedActivity = po.Activity{
		Name:        "add",
		URI:         "activity:testing_create_or_update_namespace/add",
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
		Namespace:        "testing_create_or_update_namespace",
		Comment:          "usage for workflow pass_task",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, workflow.Data.URI, "statemachine:testing_create_or_update_namespace/pass_task")
	assert.Equal(t, workflow.Data.Name, "pass_task")
	assert.Equal(t, workflow.Data.Comment, "usage for workflow pass_task")
	assert.Equal(t, workflow.Data.Status, "Enable")
	assert.Equal(t, workflow.Data.NamespaceID, ns.Data.ID)
	// 更新
	workflow, err = myTemplateService.CreateOrUpdateStateMachine(ctx, vo.CreateStateMachineRequest{
		StateMachineName: "pass_task",
		Namespace:        "testing_create_or_update_namespace",
		Comment:          "usage for workflow pass_task update",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, workflow.Data.URI, "statemachine:testing_create_or_update_namespace/pass_task")
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
