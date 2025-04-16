package template

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
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
