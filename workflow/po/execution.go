package po

import (
	"time"
)

// Execution  a workflow execution instance
type Execution struct {
	ID              int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID            string     `json:"uuid" gorm:"not null; type:VARCHAR(255);unique; comment:UUID"` //execution uuid
	URI             string     `json:"uri" gorm:"not null; index; size:255; comment:flow uri"`       //flow uri
	Status          string     `json:"status" gorm:"not null; type:VARCHAR(100)"`                    //状态
	Title           string     `json:"title" gorm:"not null; size:255"`                              // execution title
	MaxExecuteIndex int        `json:"max_execute_index" gorm:"type:INT(10);not null; default 0; comment:max execute index"`
	Definition      string     `json:"definition"  gorm:"type:MEDIUMTEXT;not null"`          // workflow  content
	Header          string     `json:"header" gorm:"not null; type:JSON"`                    // Execution 头信息
	Data            string     `json:"data" gorm:"type:MEDIUMTEXT"`                          // 执行数据
	Input           string     `json:"input" gorm:"type:MEDIUMTEXT"`                         //输入
	Output          string     `json:"output" gorm:"type:MEDIUMTEXT"`                        // 输出
	Exception       string     `json:"exception" gorm:"type:MEDIUMTEXT"`                     //异常信息
	StartTime       *time.Time `json:"start_time" gorm:"type:TIMESTAMP null; default:null"`  //开始时间
	FinishTime      *time.Time `json:"finish_time" gorm:"type:TIMESTAMP null; default:null"` // 结束时间
	ExpireTime      *time.Time `json:"expire_time" gorm:"type:TIMESTAMP null; default:null"` //超时时间
	ExecuteCount    int        `json:"execute_count" gorm:"type:INT(10);default:0"`          //执行次数
	UpdateTime      time.Time  `json:"update_time" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	CreateTime      time.Time  `json:"create_time" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}

// Step  step in workflow execution
type Step struct {
	ID          int `json:"id" gorm:"primaryKey;autoIncrement"`
	ExecutionID int `json:"execution_id"  gorm:"not null;uniqueIndex:uni_execution_index;uniqueIndex:uni_group_index; uniqueIndex:uni_group_state_name; comment:execution_id"`
	//ExecuteIndex  execute index in Execution
	ExecuteIndex int `json:"execute_index" gorm:"not null;uniqueIndex:uni_execution_index;type:INT(11);comment:ExecuteIndex"`
	//GroupID  group belongs to, 0 means no group
	GroupID      int        `json:"group_id" gorm:"not null; uniqueIndex:uni_group_state_name;uniqueIndex:uni_group_index;type:INT(11)"`
	Name         string     `json:"name" gorm:"not null;uniqueIndex:uni_group_state_name;type:VARCHAR(255)"` // State name
	GroupIndex   int        `json:"group_index" gorm:"not null;uniqueIndex:uni_group_index;"`                // State name
	Status       string     `json:"status" gorm:"not null ;type:VARCHAR(100)"`                               // State status
	Depth        int        `json:"depth" gorm:" not null; type:INT(11)"`                                    // depth in execution tree
	Type         string     `json:"type" gorm:"not null;type:VARCHAR(100)"`                                  // State type
	Resource     string     `json:"resource" gorm:"type:VARCHAR(255)"`                                       // task resource URI
	Definition   string     `json:"definition" gorm:"not null;type:MEDIUMTEXT"`                              // state definition
	Input        string     `json:"input" gorm:"type:MEDIUMTEXT;DEFAULT:null"`                               //state input
	Output       string     `json:"output" gorm:"type:MEDIUMTEXT;DEFAULT:null"`                              // state output
	Exception    string     `json:"exception" gorm:"type:MEDIUMTEXT;default:null"`                           // state exception
	Data         string     `json:"data" gorm:"type:JSON"`                                                   // state inner data
	StartTime    *time.Time `json:"start_time" gorm:"type:TIMESTAMP null; DEFAULT:null"`                     // start time
	FinishTime   *time.Time `json:"finish_time" gorm:"type:TIMESTAMP null; DEFAULT:null"`                    // finish time
	ExpireTime   *time.Time `json:"expire_time" gorm:"type:TIMESTAMP null; DEFAULT:null"`                    // expire time
	ExecuteCount int        `json:"execute_count" gorm:"not null;type:INT(10); default:0;comment:执行次数 "`     // execute count
	UpdateTime   time.Time  `json:"update_time" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	CreateTime   time.Time  `json:"create_time" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}

// StepGroup  StepGroup 是一个状态组，包含多个子状态
type StepGroup struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement"`
	// 自身的StepID， 每个组都有有个实体的State记录, 每个组类型的只有一个Stategroup
	StepID int `json:"step_id" gorm:"not null; unique;type:INT(11)"`
	// Group 所属的ExecutionID
	ExecutionID int `json:"execution_id" gorm:"not null; uniqueIndex:uni_execution_group;type:INT(11)"`
	//group所属StateID, 比如 parallel/Map 包含多个子组
	MasterStepID int `json:"master_step_id" gorm:"not null; uniqueIndex:uni_subgroup_index;type:INT(11)"`
	// GroupIndex在 MasterState 中的排序
	GroupIndex int `json:"group_index" gorm:"not null;uniqueIndex:uni_subgroup_index;type:INT(11)"`
	//SubGroupID  该StateGroup所管理group
	SubGroupID int `json:"sub_group_id" gorm:"not null; uniqueIndex:uni_execution_group;type:INT(11)"`
	// StartAt 组内开始节点的名字
	StartAt string `json:"start_at" gorm:"not null;type:VARCHAR(255)"`
	//结束节点 组内结束节点的名字
	LastAt string `json:"last_at" gorm:"not null;type:VARCHAR(255)"`

	UpdateTime time.Time `json:"update_time" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	CreateTime time.Time `json:"create_time" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}

// ExecutionEvent execution 执行日志事件

// TaskToken task 执行的token
type TaskToken struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement;type:INT(11)"`
	//关联的 uuid token uuid
	Token string `json:"token" gorm:"not null;unique;type:VARCHAR(255)"`
	// 关联的 State_id, 每个State 只能关联一个 take_token,只能有一个正在运行的任务
	StepID     int       `json:"step_id" gorm:"not null;unique;type:INT(11)"`
	IsDeleted  bool      `json:"is_deleted" gorm:"not null;type:bool"`
	UpdateTime time.Time `json:"update_time" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	CreateTime time.Time `json:"create_time" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}

// ActivityTask  staged running activity task
type ActivityTask struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement;type:BIGINT"`
	// ExecutionID
	ExecutionID int `json:"execution_id" gorm:"not null;index;type: INT(10)"` // execution ID
	//  state id
	StepID int `json:"State_id" gorm:"not null; unique;type:INT(10)"`
	// resource activity uri
	Resource string `json:"resource" gorm:"not null;index;type:VARCHAR(255)" `
	// input json as string
	Input string `json:"input" gorm:"type:JSON"` // input
	// task token ,唯一
	Token string `json:"token" gorm:"not null;unique; type:VARCHAR(255)"` // task token
	// task extra data
	Data       string    `json:"data" gorm:"type:MEDIUMTEXT"`
	CreateTime time.Time `json:"create_time" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}
