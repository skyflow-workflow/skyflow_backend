package template

import (
	"context"
	"log/slog"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mmtbak/microlibrary/paging"
	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/gorm/schema"
)

func (svc *templateService) CleanTestTableData(ctx context.Context, tx rdb.Tx, table any) error {

	var err error

	tx, maker := svc.dbClient.NewTxMaker(tx)
	defer maker.Close(&err)

	schema, err := schema.Parse(table, &sync.Map{}, svc.dbClient.DB().NamingStrategy)
	if err != nil {
		return err
	}
	tableName := schema.Table
	err = tx.Exec("TRUNCATE TABLE " + tableName).Error
	return err

}

func TestNamespaceOperation(t *testing.T) {

	var err error

	myTemplateService := NewTemplateService(getTestDBClient())
	err = myTemplateService.CleanTestTableData(context.Background(), nil, &po.Namespace{})
	assert.Equal(t, err, nil)

	nsName := "unittest"
	// create namespace
	createNsResp, err := myTemplateService.CreateNamespace(context.Background(), vo.CreateNamespaceRequest{
		Name:        nsName,
		Description: nsName,
	}, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, createNsResp.Data.Name, nsName)

	// describe namespace
	describeNsResp, err := myTemplateService.DescribeNamespace(context.Background(), nsName, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, describeNsResp.Name, nsName)

	// list namespaces
	listNsResp, err := myTemplateService.ListNamespaces(context.Background(), vo.ListNamespacesRequest{
		PageRequest: paging.DefaultPageRequest,
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, len(listNsResp.Namespaces), 1)
	assert.Equal(t, listNsResp.Namespaces[0].Name, nsName)
}

func TestActivityOperation(t *testing.T) {

	var err error
	ctx := context.Background()

	myTemplateService := NewTemplateService(getTestDBClient())
	err = myTemplateService.CleanTestTableData(context.Background(), nil, new(po.Activity))
	assert.Equal(t, err, nil)
	unittestNsName := "unittest_namespace"
	unittestNsDescription := "unittest_namespace_description"

	// create namespace
	createNsResp, err := myTemplateService.CreateNamespace(ctx, vo.CreateNamespaceRequest{
		Name:        unittestNsName,
		Description: unittestNsDescription,
	}, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, createNsResp.Data.Name, unittestNsName)
	assert.Equal(t, createNsResp.Data.Description, unittestNsDescription)

	// create activity
	unittestActivity := po.Activity{
		Name:        "unitest_activity_name",
		Description: "unitest_activity_description",
		Parameters:  `{"a": "int", "b": "int"}`,
	}

	createActivityResp, err := myTemplateService.CreateActivity(ctx, vo.CreateActivityRequest{
		ActivityName: unittestActivity.Name,
		Namespace:    unittestNsName,
		Description:  unittestActivity.Description,
		Parameters:   unittestActivity.Parameters,
	}, nil)
	assert.Equal(t, err, nil)

	expectedActivity := po.Activity{
		Name:        "unitest_activity_name",
		URI:         "activity:unittest_namespace/unitest_activity_name",
		Description: "unitest_activity_description",
		Parameters:  `{"a": "int", "b": "int"}`,
		Status:      "Enable",
		NamespaceID: createNsResp.Data.ID,
	}
	assert.Equal(t, createActivityResp.Data.Parameters, expectedActivity.Parameters)
	assert.Equal(t, createActivityResp.Data.URI, expectedActivity.URI)
	assert.Equal(t, createActivityResp.Data.Name, expectedActivity.Name)
	assert.Equal(t, createActivityResp.Data.Description, expectedActivity.Description)
	assert.Equal(t, createActivityResp.Data.Status, expectedActivity.Status)
	assert.Equal(t, createActivityResp.Data.NamespaceID, expectedActivity.NamespaceID)

	// describe activity
	describeActivityResp, err := myTemplateService.DescribeActivity(ctx, expectedActivity.URI, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, describeActivityResp.URI, expectedActivity.URI)
	assert.Equal(t, describeActivityResp.Parameters, expectedActivity.Parameters)
	assert.Equal(t, describeActivityResp.Name, expectedActivity.Name)
	assert.Equal(t, describeActivityResp.Description, expectedActivity.Description)

	// list activities
	listActivitiesResp, err := myTemplateService.ListActivities(ctx, vo.ListActivitiesRequest{
		PageRequest: paging.DefaultPageRequest,
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, len(listActivitiesResp.Activities), 1)
	assert.Equal(t, listActivitiesResp.Activities[0].Name, unittestActivity.Name)
	assert.Equal(t, listActivitiesResp.Activities[0].Description, unittestActivity.Description)
	// delete activity
	err = myTemplateService.DeleteActivity(ctx, vo.DeleteActivityRequest{
		ActivityURI: expectedActivity.URI,
	}, nil)
	assert.Equal(t, err, nil)
	// list activities after delete
	listActivitiesResp, err = myTemplateService.ListActivities(ctx, vo.ListActivitiesRequest{
		PageRequest: paging.DefaultPageRequest,
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, len(listActivitiesResp.Activities), 0)
	// describe activity after delete
	_, err = myTemplateService.DescribeActivity(ctx, expectedActivity.URI, nil)
	assert.Equal(t, err.Error(), "record not found")

}

func TestStateMachineOperation(t *testing.T) {

	var err error
	ctx := context.Background()

	myTemplateService := NewTemplateService(getTestDBClient())
	err = myTemplateService.CleanTestTableData(context.Background(), nil, &po.StateMachine{})
	assert.Equal(t, err, nil)

	unittestNamespace := po.Namespace{
		Name:        "unittest_namespace",
		Description: "unittest_namespace_description",
	}

	createNsResp, err := myTemplateService.CreateOrUpdateNamespace(ctx, vo.CreateNamespaceRequest{
		Name:        unittestNamespace.Name,
		Description: unittestNamespace.Description,
	}, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, createNsResp.Data.Name, unittestNamespace.Name)
	assert.Equal(t, createNsResp.Data.Description, unittestNamespace.Description)

	slog.Info("Created or updated namespace",
		"namespace", unittestNamespace.Name,
		"description", unittestNamespace.Description,
		"namespaceID", createNsResp.Data.ID)

	// create state machine
	unittestStateMachine := po.StateMachine{
		Name:        "unittest_state_machine",
		Description: "unittest_state_machine_description",
		Definition:  `{"StartAt": "Pass", "States": {"Pass": {"Type": "Pass", "End": true}}}`,
	}
	createStateMachineResp, err := myTemplateService.CreateStateMachine(ctx, vo.CreateStateMachineRequest{
		Name:        unittestStateMachine.Name,
		Namespace:   unittestNamespace.Name,
		Description: unittestStateMachine.Description,
		Definition:  unittestStateMachine.Definition,
	}, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, createStateMachineResp.Data.URI, "statemachine:unittest_namespace/unittest_state_machine")
	assert.Equal(t, createStateMachineResp.Data.Name, unittestStateMachine.Name)
	assert.Equal(t, createStateMachineResp.Data.Description, unittestStateMachine.Description)
	assert.Equal(t, createStateMachineResp.Data.Status, "Enable")
	assert.Equal(t, createStateMachineResp.Data.NamespaceID, createNsResp.Data.ID)

	unittestStateMachine.Description = "unittest_namespace_description_updated"
	unittestStateMachine.Definition = `{"StartAt": "Pass", "States": {"Pass": {"Type": "Pass", "End": true, "Result": "Hello World"}}}`
	// update state machine
	createStateMachineResp, err = myTemplateService.CreateOrUpdateStateMachine(ctx, vo.CreateStateMachineRequest{
		Name:        unittestStateMachine.Name,
		Namespace:   unittestNamespace.Name,
		Description: unittestStateMachine.Description,
		Definition:  unittestStateMachine.Definition,
	}, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, createStateMachineResp.Data.URI, "statemachine:unittest_namespace/unittest_state_machine")
	assert.Equal(t, createStateMachineResp.Data.Name, unittestStateMachine.Name)
	assert.Equal(t, createStateMachineResp.Data.Description, unittestStateMachine.Description)
	assert.Equal(t, createStateMachineResp.Data.Definition, unittestStateMachine.Definition)
	assert.Equal(t, createStateMachineResp.Data.Status, "Enable")
	assert.Equal(t, createStateMachineResp.Data.NamespaceID, createNsResp.Data.ID)

	// describe state machine
	describeStateMachineResp, err := myTemplateService.DescribeStateMachine(ctx, vo.DescribeStateMachineRequest{
		StateMachineURI: createStateMachineResp.Data.URI,
	}, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, describeStateMachineResp.URI, createStateMachineResp.Data.URI)
	assert.Equal(t, describeStateMachineResp.Name, createStateMachineResp.Data.Name)
	assert.Equal(t, describeStateMachineResp.Description, unittestStateMachine.Description)
	assert.Equal(t, describeStateMachineResp.Definition, unittestStateMachine.Definition)
	assert.Equal(t, describeStateMachineResp.Status, "Enable")
	assert.Equal(t, describeStateMachineResp.NamespaceID, createNsResp.Data.ID)
	// list state machines
	listStateMachinesResp, err := myTemplateService.ListStateMachines(ctx, vo.ListStateMachinesRequest{
		PageRequest: paging.DefaultPageRequest,
	}, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(listStateMachinesResp.StateMachines), 1)
	assert.Equal(t, listStateMachinesResp.StateMachines[0].Name, unittestStateMachine.Name)
	assert.Equal(t, listStateMachinesResp.StateMachines[0].Description, unittestStateMachine.Description)
	assert.Equal(t, listStateMachinesResp.StateMachines[0].Definition, unittestStateMachine.Definition)
	assert.Equal(t, listStateMachinesResp.StateMachines[0].URI, createStateMachineResp.Data.URI)
	// delete state machine
	err = myTemplateService.DeleteStateMachine(ctx, vo.DeleteStateMachineRequest{
		StateMachineURI: createStateMachineResp.Data.URI,
	}, nil)
	assert.Equal(t, err, nil)
	// list state machines after delete
	listStateMachinesResp, err = myTemplateService.ListStateMachines(ctx, vo.ListStateMachinesRequest{
		PageRequest: paging.DefaultPageRequest,
	}, nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(listStateMachinesResp.StateMachines), 0)
	// describe state machine after delete
	_, err = myTemplateService.DescribeStateMachine(ctx, vo.DescribeStateMachineRequest{
		StateMachineURI: createStateMachineResp.Data.URI,
	}, nil)
	assert.Equal(t, err.Error(), "record not found")
}
