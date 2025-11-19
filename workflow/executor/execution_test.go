package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"

	"gopkg.in/go-playground/assert.v1"
)

func TestSort(t *testing.T) {

	type Family struct {
		Name string
		Age  int
	}

	family1 := []Family{
		{"Alice", 23},
		{"David", 2},
		{"Eve", 2},
		{"Bob", 25},
	}
	family2 := []Family{
		{"Alice", 23},
		{"David", 2},
		{"Eve", 2},
		{"Bob", 25},
	}
	family3 := []Family{
		{"Alice", 14},
		{"David", 2},
		{"Eve", 2},
		{"Bob", 25},
	}

	families := map[string][]Family{
		"A": family1,
		"B": family2,
		"C": family3,
	}

	for _, family := range families {
		// Sort by age, keeping original order or equal elements.
		sort.SliceStable(family, func(i, j int) bool {
			return family[i].Age < family[j].Age
		})
	}
	for n, family := range families {
		fmt.Println(n, "-------", family)
	}
}

func TestExecution(t *testing.T) {
	var err error
	var dbExe *po.Execution
	var executionID = 752700

	dbExe, err = myExecutionService.QueryExecutionByID(executionID, ExecutionFields.L5, nil)
	if err != nil {
		t.Error(err)
		return
	}

	exe, err := NewExecutionFromData(dbExe, myExecutionService)
	fmt.Println(err)
	fmt.Println(exe)
	// bone := se.GetBone()
	// fmt.Println("Get bone Success")
	// // fmt.Println(bone)
	// bonbyte, err := json.Marshal(bone)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(bonbyte))

}

func TestCreateExecution(t *testing.T) {

	templatepath := "../tests/templates"
	fileInfo, err := os.Stat(templatepath)
	assert.Equal(t, err, nil)
	assert.Equal(t, fileInfo.IsDir(), true)

	fileInfos, err := os.ReadDir(templatepath)
	assert.Equal(t, err, nil)

	var testcases = []struct {
		template string
		input    string
	}{{
		template: `
		{
			"StartAt":"PStart",
			"States":{
				"PStart":{
					"Type":"Pass",
					"End":true
				}
			}
		}`,
		input: `{}`,
	},
	}

	for _, info := range fileInfos {
		if strings.HasSuffix(info.Name(), ".json") {
			fpath := filepath.Join(templatepath, info.Name())
			content, err := os.ReadFile(fpath)
			assert.Equal(t, err, nil)
			testcases = append(testcases,
				struct {
					template string
					input    string
				}{
					template: string(content),
					input:    `{}`,
				})
		}
	}

	for idx, tt := range testcases {
		fmt.Println("index  :", idx)
		req := vo.StartExecutionRequest{
			StateMachineDefinition: tt.template,
			Input:                  tt.input,
		}
		// 创建
		dbExe, err := myExecutionService.StartExecution(req)
		fmt.Println(err)
		assert.Equal(t, err == nil, true)
		fmt.Println(dbExe)
		// 初始化
		msg := queue.InnerMessageBody{
			ExecutionID: dbExe.ID,
			Type:        MessageType.ExecutionInit,
		}

		exe, err := NewExecutionFromData(dbExe, myExecutionService)
		fmt.Println(err)
		assert.Equal(t, err == nil, true)
		err = exe.ProcessEvent(msg)
		fmt.Println(err)

	}
}

func TestExecutionProcessMessage(t *testing.T) {

	var testcases = []struct {
		template string
		input    string
	}{
		{
			template: `
			{
				"StartAt":"P1",
				"States":{
					"P1":{
						"Type":"Pass",
						"End":true
					}
				}
			}`,
			input: "{}",
		},
	}

	for idx, tt := range testcases {

		fmt.Println("index : ", idx)
		// var dbexecution db.Execution
		req := vo.StartExecutionRequest{
			Title:                  "unit test execution ",
			StateMachineDefinition: tt.template,
			Input:                  tt.input,
		}
		dbExe, err := myExecutionService.StartExecution(req)
		fmt.Println(err)
		assert.Equal(t, err == nil, true)

		execution_id := dbExe.ID
		fmt.Println("execution id :", execution_id)

		// exe, err := NewExecutionFromID(dbexe.ID)
		exe, err := NewExecutionFromData(dbExe, myExecutionService)
		fmt.Println(err)
		assert.Equal(t, err == nil, true)
		msg := queue.InnerMessageBody{
			ExecutionID: execution_id,
			Type:        MessageType.ExecutionInit,
		}
		err = exe.ProcessEvent(msg)
		fmt.Println(err)
		assert.Equal(t, err == nil, true)

		msg = queue.InnerMessageBody{
			ExecutionID: execution_id,
			Type:        MessageType.ExecutionSucceed,
		}
		err = exe.ProcessEvent(msg)
		fmt.Println(err)
		assert.Equal(t, err == nil, true)
	}

}

// func TestExecutionBone(t *testing.T) {

// 	var testcases = []struct {
// 		uuid  string
// 		input string
// 	}{
// 		{
// 			uuid: "8d017491-b282-11ec-abd6-52540002d945",
// 		},
// 		{
// 			uuid: "47dde0cf-af50-11ec-8fa6-52540002d945",
// 		},
// 		{
// 			uuid: "f58feb2d-b1b0-11ec-a83b-52540002d945",
// 		},
// 	}
// 	for _, tt := range testcases {
// 		exe, err := myExecutionService.NewExecutionFromUUID(tt.uuid, ExecutionFields.L3, nil)
// 		assert.Equal(t, err, nil)
// 		bone, err := exe.GetBone()
// 		assert.Equal(t, err, nil)
// 		fmt.Println(bone)

// 	}
// }

// func TestProcessExecutionInit(t *testing.T) {

// 	uuid := "6b34e12f-e935-4dd0-9a69-b38a4b7d8031"

// 	exe, err := myExecutionService.NewExecutionFromUUID(uuid, ExecutionFields.L1, nil)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = exe.ProcessInit()
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// }
