package states

// TaskBody ...
type TaskBody struct {
	Resource         string `validate:"required,gt=0"`
	TimeoutSeconds   uint   `validate:"gte=0"`
	HeartbeatSeconds uint   `validate:"gte=0"`
	// Retry for decode
	Retry []RetryNode ``
	// Catch for decode
	Catch []CatchNode ``
}

// RetryNode struct for retry
type RetryNode struct {
	ErrorEquals     []string `mapstructure:"ErrorEquals"`
	IntervalSeconds uint     `mapstructure:"IntervalSeconds"`
	MaxAttempts     uint     `mapstructure:"MaxAttempts"`
	BackoffRate     float64  `mapstructure:"BackoffRate"`
}

// CatchNode  struct for catch
type CatchNode struct {
	ErrorEquals []string `mapstructure:"ErrorEquals"`
	Next        string   `mapstructure:"Next"`
	ResultPath  string   `mapstructure:"ResultPath"`
}

// TaskSendData  Task State执行结果发送的数据
type TaskSendData struct {
	Success bool        // 执行是否成功
	Retry   []int       // 当前的Retry次数
	Errors  []string    // 发送的Error列表
	Output  interface{} // Task 节点的执行输出
}

// DefaultRetryNode default retry config for task
var DefaultRetryNode = RetryNode{
	ErrorEquals:     []string{},
	IntervalSeconds: 1,
	MaxAttempts:     3,
	BackoffRate:     1.5,
}

// DefaultCatchNode default catch config for task
var DefaultCatchNode = CatchNode{
	ErrorEquals: []string{},
	Next:        "",
	ResultPath:  "$",
}

// DefaultTaskBody ...
var DefaultTaskBody = TaskBody{
	Resource: "",
	// 默认0, 0 即没有超时限制
	TimeoutSeconds:   0,
	HeartbeatSeconds: 0, // 默认0, 0 即没有超时限制
}

// Validate ...
func (body *TaskBody) Validate() error {
	var err error
	err = myvalidate.Struct(body)
	if err != nil {
		return err
	}
	_, err = ParseResource(body.Resource)
	if err != nil {
		return err
	}
	return err
}

// Init ...
func (body *TaskBody) Init() error {
	return body.Validate()
}

// Task ...
type Task struct {
	*BaseState
	*TaskBody
}

// Init ...
func (t *Task) Init() error {
	// Catch
	if len(t.TaskBody.Catch) > 0 {
		// retry 不为空
		var nexts = []string{}
		for _, catchnode := range t.TaskBody.Catch {
			nexts = append(nexts, catchnode.Next)
		}
		// 添加Next 到bone的next 中
		t._Bone.Next = append(t._Bone.Next, nexts...)
	}
	return nil
}
