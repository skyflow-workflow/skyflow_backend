package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

// MockExecutionService 模拟ExecutionService
type MockExecutionService struct {
	mock.Mock
}

func (m *MockExecutionService) QueryExecutionByID(id int, fields []string, tx interface{}) (*po.Execution, error) {
	args := m.Called(id, fields, tx)
	return args.Get(0).(*po.Execution), args.Error(1)
}

func (m *MockExecutionService) QueryExecutionURL(uuid string) string {
	args := m.Called(uuid)
	return args.String(0)
}

func (m *MockExecutionService) NewExecutionFromUUID(uuid string, fields []string, tx interface{}) (*Execution, error) {
	args := m.Called(uuid, fields, tx)
	return args.Get(0).(*Execution), args.Error(1)
}

func (m *MockExecutionService) StartExecution(req vo.StartExecutionRequest) (*po.Execution, error) {
	args := m.Called(req)
	return args.Get(0).(*po.Execution), args.Error(1)
}

func (m *MockExecutionService) DomainService() interface{} {
	return nil
}

func (m *MockExecutionService) LockService() interface{} {
	return nil
}

func (m *MockExecutionService) MetaDB() interface{} {
	return nil
}

func (m *MockExecutionService) InnerQueue() interface{} {
	return nil
}

func (m *MockExecutionService) StandardExecutor() interface{} {
	return nil
}

func (m *MockExecutionService) SendExecutionEvents(events ...vo.ExecutionEvent) {
	m.Called(events)
}

func (m *MockExecutionService) ExecutionErrorProcess(err error, executionID int) error {
	args := m.Called(err, executionID)
	return args.Error(0)
}

// MockStateMachine 模拟StateMachine
type MockStateMachine struct {
	mock.Mock
}

func (m *MockStateMachine) GetInput(input interface{}, info ExecutionInfo) (interface{}, error) {
	args := m.Called(input, info)
	return args.Get(0), args.Error(1)
}

func (m *MockStateMachine) GetTimeout() states.Timeout {
	args := m.Called()
	return args.Get(0).(states.Timeout)
}

func (m *MockStateMachine) StateMachineHeader() interface{} {
	return nil
}

func (m *MockStateMachine) StartAt() string {
	return ""
}

func init() {
}

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

	familys := map[string][]Family{
		"A": family1,
		"B": family2,
		"C": family3,
	}

	for _, family := range familys {
		// Sort by age, keeping original order or equal elements.
		sort.SliceStable(family, func(i, j int) bool {
			return family[i].Age < family[j].Age
		})
	}
	for n, family := range familys {
		fmt.Println(n, "-------", family)
	}
}

func TestExecution(t *testing.T) {
	var err error
	var dbexe po.Execution
	var executionID = 752700
	// var dbsteps []po.Step
	// var dbstepgroup []po.StepGroup
	// db.GetModel().GetEngine().Id(executionID).Get(&dbsm)
	// db.GetModel().GetEngine().Where("execution_id = ?", executionID).Find(&dbsteps)
	// db.GetModel().GetEngine().Where("execution_id = ?", executionID).Find(&dbstepgroup)

	dbexe, err = myExecutionService.QueryExecutionByID(executionID, ExecutionFields.L5, nil)
	if err != nil {
		t.Error(err)
		return
	}

	exe, err := NewExecutionFromData(&dbexe, myExecutionService)
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
	finfo, err := os.Stat(templatepath)
	assert.Equal(t, err, nil)
	assert.Equal(t, finfo.IsDir(), true)

	finfos, err := os.ReadDir(templatepath)
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

	for _, info := range finfos {
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
			WorkflowDefinition: tt.template,
			Input:              tt.input,
		}
		// 创建
		dbexe, err := myExecutionService.StartExecution(req)
		fmt.Println(err)
		assert.Equal(t, err == nil, true)
		fmt.Println(dbexe)
		// 初始化
		msg := queue.InnerMessage{
			ExecutionID: dbexe.ID,
			Type:        MessageType.ExecutionInit,
		}

		exe, err := NewExecutionFromData(&dbexe, myExecutionService)
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
			Title:              "unit test execution ",
			WorkflowDefinition: tt.template,
			Input:              tt.input,
		}
		dbexe, err := myExecutionService.StartExecution(req)
		fmt.Println(err)
		assert.Equal(t, err == nil, true)

		execution_id := dbexe.ID
		fmt.Println("execution id :", execution_id)

		// exe, err := NewExecutionFromID(dbexe.ID)
		exe, err := NewExecutionFromData(&dbexe, myExecutionService)
		fmt.Println(err)
		assert.Equal(t, err == nil, true)
		msg := queue.InnerMessage{
			ExecutionID: execution_id,
			Type:        MessageType.ExecutionInit,
		}
		err = exe.ProcessEvent(msg)
		fmt.Println(err)
		assert.Equal(t, err == nil, true)

		msg = queue.InnerMessage{
			ExecutionID: execution_id,
			Type:        MessageType.ExecutionSucceed,
		}
		err = exe.ProcessEvent(msg)
		fmt.Println(err)
		assert.Equal(t, err == nil, true)
	}

}

func TestExecutionBone(t *testing.T) {

	var testcases = []struct {
		uuid  string
		input string
	}{
		{
			uuid: "8d017491-b282-11ec-abd6-52540002d945",
		},
		{
			uuid: "47dde0cf-af50-11ec-8fa6-52540002d945",
		},
		{
			uuid: "f58feb2d-b1b0-11ec-a83b-52540002d945",
		},
	}
	for _, tt := range testcases {
		exe, err := myExecutionService.NewExecutionFromUUID(tt.uuid, ExecutionFields.L3, nil)
		assert.Equal(t, err, nil)
		bone, err := exe.GetBone()
		assert.Equal(t, err, nil)
		fmt.Println(bone)

	}
}

func TestProcessExecutionInit(t *testing.T) {

	uuid := "6b34e12f-e935-4dd0-9a69-b38a4b7d8031"

	exe, err := myExecutionService.NewExecutionFromUUID(uuid, ExecutionFields.L1, nil)
	if err != nil {
		t.Error(err)
		return
	}
	err = exe.ProcessInit()
	if err != nil {
		t.Error(err)
		return
	}
}

// TestGetInput_NormalCase 测试GetInput方法的正常情况
func TestGetInput_NormalCase(t *testing.T) {
	// 准备测试数据
	mockService := new(MockExecutionService)
	mockStateMachine := new(MockStateMachine)

	// 创建Execution实例
	execution := &Execution{
		Data: &po.Execution{
			UUID:  "test-uuid-123",
			Input: `{"name": "test", "value": 123}`,
		},
		StateMachine:     mockStateMachine,
		ExecutionService: mockService,
	}

	// 设置模拟行为
	mockService.On("QueryExecutionURL", "test-uuid-123").Return("http://test.com/executions/test-uuid-123")
	mockStateMachine.On("GetInput", mock.Anything, mock.Anything).Return(map[string]interface{}{
		"name":  "test",
		"value": 123,
		"execution": ExecutionInfo{
			UUID: "test-uuid-123",
			URL:  "http://test.com/executions/test-uuid-123",
		},
	}, nil)

	// 执行测试
	result, err := execution.GetInput()

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)

	resultMap, ok := result.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "test", resultMap["name"])
	assert.Equal(t, 123, resultMap["value"])

	// 验证模拟调用
	mockService.AssertCalled(t, "QueryExecutionURL", "test-uuid-123")
	mockStateMachine.AssertCalled(t, "GetInput", mock.Anything, mock.Anything)
}

// TestGetInput_InvalidJSON 测试GetInput方法处理无效JSON的情况
func TestGetInput_InvalidJSON(t *testing.T) {
	// 准备测试数据
	execution := &Execution{
		Data: &po.Execution{
			UUID:  "test-uuid-123",
			Input: `invalid json`,
		},
		ExecutionService: new(MockExecutionService),
	}

	// 执行测试
	result, err := execution.GetInput()

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid character")
}

// TestGetInput_EmptyInput 测试GetInput方法处理空输入的情况
func TestGetInput_EmptyInput(t *testing.T) {
	// 准备测试数据
	mockService := new(MockExecutionService)
	mockStateMachine := new(MockStateMachine)

	execution := &Execution{
		Data: &po.Execution{
			UUID:  "test-uuid-123",
			Input: "",
		},
		StateMachine:     mockStateMachine,
		ExecutionService: mockService,
	}

	// 设置模拟行为
	mockService.On("QueryExecutionURL", "test-uuid-123").Return("http://test.com/executions/test-uuid-123")
	mockStateMachine.On("GetInput", mock.Anything, mock.Anything).Return(nil, nil)

	// 执行测试
	result, err := execution.GetInput()

	// 验证结果
	assert.NoError(t, err)
	assert.Nil(t, result)
}

// TestNewExecutionFromData_NormalCase 测试NewExecutionFromData方法的正常情况
func TestNewExecutionFromData_NormalCase(t *testing.T) {
	// 准备测试数据
	mockService := new(MockExecutionService)
	testData := &po.Execution{
		ID:     123,
		UUID:   "test-uuid-456",
		Status: "created",
	}

	// 执行测试
	execution, err := NewExecutionFromData(testData, mockService)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, execution)
	assert.Equal(t, testData, execution.Data)
	assert.Equal(t, mockService, execution.ExecutionService)
	assert.NotNil(t, execution.States)
	assert.Empty(t, execution.States)
}

// TestNewExecutionFromData_NilData 测试NewExecutionFromData方法处理nil数据的情况
func TestNewExecutionFromData_NilData(t *testing.T) {
	// 准备测试数据
	mockService := new(MockExecutionService)

	// 执行测试
	execution, err := NewExecutionFromData(nil, mockService)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, execution)
	assert.Nil(t, execution.Data)
	assert.Equal(t, mockService, execution.ExecutionService)
}

// TestFullInit_AlreadyInitialized 测试FullInit方法已经初始化的情况
func TestFullInit_AlreadyInitialized(t *testing.T) {
	// 准备测试数据
	mockService := new(MockExecutionService)
	mockStateMachine := new(MockStateMachine)

	execution := &Execution{
		Data: &po.Execution{
			ID: 123,
		},
		StateMachine:     mockStateMachine,
		ExecutionService: mockService,
	}

	// 执行测试
	err := execution.FullInit()

	// 验证结果
	assert.NoError(t, err)
}

// TestProcessEvent_UnrecognizedEvent 测试ProcessEvent方法处理未知事件类型
func TestProcessEvent_UnrecognizedEvent(t *testing.T) {
	// 准备测试数据
	mockService := new(MockExecutionService)

	execution := &Execution{
		Data: &po.Execution{
			ID: 123,
		},
		ExecutionService: mockService,
	}

	// 设置模拟行为
	mockService.On("LockService").Return(nil)
	mockService.On("MetaDB").Return(nil)

	// 执行测试
	msg := queue.InnerMessageBody{
		ExecutionID: 123,
		Type:        "UnknownEventType",
	}
	err := execution.ProcessEvent(msg)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unrecognized event")
}

// TestChangeExecutionStatus_NormalCase 测试ChangeExecutionStatus方法的正常情况
func TestChangeExecutionStatus_NormalCase(t *testing.T) {
	// 准备测试数据
	mockService := new(MockExecutionService)

	execution := &Execution{
		Data: &po.Execution{
			ID: 123,
		},
		ExecutionService: mockService,
	}

	// 设置模拟行为
	mockService.On("MetaDB").Return(nil)

	// 执行测试
	err := execution.ChangeExecutionStatus(ExecutionStatus.Success, nil)

	// 验证结果
	assert.NoError(t, err)
}

// TestProcessSucceed_NormalCase 测试ProcessSucceed方法的正常情况
func TestProcessSucceed_NormalCase(t *testing.T) {
	// 准备测试数据
	mockService := new(MockExecutionService)

	execution := &Execution{
		Data: &po.Execution{
			ID: 123,
		},
		ExecutionService: mockService,
	}

	// 设置模拟行为
	mockService.On("MetaDB").Return(nil)
	mockService.On("QueryExecutionByID", 123, []string{"id", "output"}, mock.Anything).Return(&po.Execution{
		ID:     123,
		Output: "{}",
	}, nil)
	mockService.On("SendExecutionEvents", mock.Anything)
	mockService.On("InnerQueue").Return(nil)

	// 执行测试
	err := execution.ProcessSucceed(`{"result": "success"}`)

	// 验证结果
	assert.NoError(t, err)
}

// TestStopExecution_NormalCase 测试StopExecution方法的正常情况
func TestStopExecution_NormalCase(t *testing.T) {
	// 准备测试数据
	mockService := new(MockExecutionService)

	execution := &Execution{
		Data: &po.Execution{
			ID:     123,
			UUID:   "test-uuid-789",
			Status: string(ExecutionStatus.Running),
		},
		ExecutionService: mockService,
	}

	// 设置模拟行为
	mockService.On("LockService").Return(nil)
	mockService.On("MetaDB").Return(nil)
	mockService.On("SendExecutionEvents", mock.Anything)
	mockService.On("InnerQueue").Return(nil)

	// 执行测试
	err := execution.StopExecution("TestError", "Test cause")

	// 验证结果
	assert.NoError(t, err)
}

// TestStopExecution_AlreadyStopped 测试StopExecution方法处理已经停止的情况
func TestStopExecution_AlreadyStopped(t *testing.T) {
	// 准备测试数据
	mockService := new(MockExecutionService)

	execution := &Execution{
		Data: &po.Execution{
			ID:     123,
			UUID:   "test-uuid-789",
			Status: string(ExecutionStatus.Aborted),
		},
		ExecutionService: mockService,
	}

	// 设置模拟行为
	mockService.On("LockService").Return(nil)
	mockService.On("MetaDB").Return(nil)

	// 执行测试
	err := execution.StopExecution("TestError", "Test cause")

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "has been stopped")
}
