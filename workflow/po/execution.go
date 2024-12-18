package po

import (
	"time"
)

// Execution  一个workflow的一个执行实例
type Execution struct {
	ID                 int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID               string     `json:"uuid" gorm:"not null; type:VARCHAR(255);unique; comment:UUID"` //execution uuid
	URI                string     `json:"uri" gorm:"not null; index; size:255; comment:flow uri"`       //flow uri
	WorkflowFlowType   string     `json:"workflow_type" gorm:"not null; size:255"`                      //workflow type
	Status             string     `json:"status" gorm:"not null; type:VARCHAR(100)"`                    //状态
	Title              string     `json:"title" gorm:"not null; size:255"`                              // execution title
	MaxExecuteIndex    int        `json:"max_execute_index" gorm:"type:INT(10);not null; default 0; comment:max execute index"`
	WorkflowDefinition string     `json:"workflow_definition"  gorm:"type:MEDIUMTEXT;not null"` // workflow  content
	Header             string     `json:"header" gorm:"not null; type:JSON"`                    // Execution 头信息
	Data               string     `json:"data" gorm:"type:MEDIUMTEXT"`                          // 执行数据
	Input              string     `json:"input" gorm:"type:MEDIUMTEXT"`                         //输入
	Output             string     `json:"ouput" gorm:"type:MEDIUMTEXT"`                         // 输出
	Exception          string     `json:"exception" gorm:"type:MEDIUMTEXT"`                     //异常信息
	StartTime          *time.Time `json:"start_time" gorm:"type:TIMESTAMP null; default:null"`  //开始时间
	FinishTime         *time.Time `json:"finish_time" gorm:"type:TIMESTAMP null; default:null"` // 结束时间
	ExpireTime         *time.Time `json:"expire_time" gorm:"type:TIMESTAMP null; default:null"` //超时时间
	ExecuteCount       int        `json:"execute_count" gorm:"type:INT(10);default:0"`          //执行次数
	GmtModified        time.Time  `json:"gmt_modified" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	GmtCreated         time.Time  `json:"gmt_created" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}

// ExecutionShade shade table for execution , for lock execution
type ExecutionShade struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID string `json:"uuid" gorm:"not null; type:VARCHAR(255);unique; comment:UUID"` //execution uuid
}

// State  一个具体 State 实例
type State struct {
	ID          int `json:"id" gorm:"primaryKey;autoIncrement"`
	ExecutionID int `json:"execution_id"  gorm:"not null;uniqueIndex:uni_execution_index;uniqueIndex:uni_group_index; uniqueIndex:uni_group_state_name; comment:execution_id"`
	//ExecuteIndex 在Execution内执行索引
	ExecuteIndex int `json:"execute_index" gorm:"not null;uniqueIndex:uni_execution_index;type:INT(11);comment:第几个执行的步骤"`
	//在Execution内的相对GroupID
	GroupID    int    `json:"group_id" gorm:"not null; uniqueIndex:uni_group_state_name;uniqueIndex:uni_group_index;type:INT(11)"`
	Name       string `json:"name" gorm:"not null;uniqueIndex:uni_group_state_name;type:VARCHAR(255)"` // State name
	GroupIndex int    `json:"group_index" gorm:"not null;uniqueIndex:uni_group_index;"`                // State name
	// status
	Status string `json:"status" gorm:"not null ;type:VARCHAR(100)"` //状态
	// State 在execution 中的深度层次
	Depth        int        `json:"depth" gorm:" not null; type:INT(11)"`
	Type         string     `json:"type" gorm:"not null;type:VARCHAR(100)"`                              // State type
	ActivityURI  string     `json:"activity_uri" gorm:"type:VARCHAR(255)"`                               //关联的activity uri
	Definition   string     `json:"definition" gorm:"not null;type:MEDIUMTEXT"`                          // State content
	Input        string     `json:"input" gorm:"type:MEDIUMTEXT;DEFAULT:null"`                           //输入
	Output       string     `json:"ouput" gorm:"type:MEDIUMTEXT;DEFAULT:null"`                           // 输出
	Exception    string     `json:"exception" gorm:"type:MEDIUMTEXT;default:null"`                       //异常信息
	Data         string     `json:"data" gorm:"type:JSON"`                                               // 内部记录state 状态
	StartTime    *time.Time `json:"start_time" gorm:"type:TIMESTAMP null; DEFAULT:null"`                 //开始时间
	FinishTime   *time.Time `json:"finish_time" gorm:"type:TIMESTAMP null; DEFAULT:null"`                // 结束时间
	ExpireTime   *time.Time `json:"expire_time" gorm:"type:TIMESTAMP null; DEFAULT:null"`                //超时时间
	ExecuteCount int        `json:"execute_count" gorm:"not null;type:INT(10); default:0;comment:执行次数 "` //执行次数
	GmtModified  time.Time  `json:"gmt_modified" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	GmtCreated   time.Time  `json:"gmt_created" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}

// StateGroup , 存储State归属Group组之间的关系， 仅仅表示归属关系
type StateGroup struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement"`
	// 自身的StateID， 每个组都有有个实体的State记录, 每个组类型的只有一个Stategroup
	StateID int `json:"State_id" gorm:"not null; unique;type:INT(11)"`
	// Group 所属的ExecutionID
	ExecutionID int `json:"execution_id" gorm:"not null; uniqueIndex:uni_execution_group;type:INT(11)"`
	//group所属StateID, 比如 parallel/Map 包含多个子组
	MasterStateID int `json:"master_state_id" gorm:"not null; uniqueIndex:uni_subgroup_index;type:INT(11)"`
	// GroupIndex在 MasterState 中的排序
	GroupIndex int `json:"group_index" gorm:"not null;uniqueIndex:uni_subgroup_index;type:INT(11)"`
	//SubGroupID  该StateGroup所管理group
	SubGroupID int `json:"sub_group_id" gorm:"not null; uniqueIndex:uni_execution_group;type:INT(11)"`
	// StartAt 组内开始节点的名字
	StartAt string `json:"start_at" gorm:"not null;type:VARCHAR(255)"`
	//结束节点 组内结束节点的名字
	LastAt string `json:"last_at" gorm:"not null;type:VARCHAR(255)"`

	GmtModified time.Time `json:"gmt_modified" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	GmtCreated  time.Time `json:"gmt_created" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}

// ExecutionEvent execution 执行日志事件

// TaskToken task 执行的token
type TaskToken struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement;type:INT(11)"`
	//关联的 uuid token uuid
	Token string `json:"token" gorm:"not null;unique;type:VARCHAR(255)"`
	// 关联的 State_id, 每个State 只能关联一个 take_token,只能有一个正在运行的任务
	StateID   int  `json:"State_id" gorm:"not null;unique;type:INT(11)"`
	IsDeleted bool `json:"is_deleted" gorm:"not null;type:bool"`
	// GmtModified time.Time `json:"gmt_modified" gorm:"autoUpdateTime;type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	// GmtCreated  time.Time `json:"gmt_created" gorm:"autoCreateTime;type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP"`
	GmtModified time.Time `json:"gmt_modified" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	GmtCreated  time.Time `json:"gmt_created" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}

// ActivityTask  暂存活动使用的数据结构
type ActivityTask struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement;type:BIGINT"`
	// ExecutionID
	ExecutionID int `json:"execution_id" gorm:"not null;index;type: INT(10)"` // execution ID
	//  state id
	StateID int `json:"State_id" gorm:"not null; unique;type:INT(10)"`
	// resource activityuri
	Resource string `json:"resource" gorm:"not null;index;type:VARCHAR(255)" `
	// input json as string
	Input string `json:"input" gorm:"type:JSON"` // input
	// task token ,唯一
	Token string `json:"token" gorm:"not null;unique; type:VARCHAR(255)"` // task token
	// task extra data
	Data string `json:"data" gorm:"type:MEDIUMTEXT"`
	// GmtCreated time.Time `json:"gmt_created" gorm:"autoCreateTime;type:TIMESTAMP;DEFAULT:CURRENT_TIMESTAMP"`
	GmtCreated time.Time `json:"gmt_created" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}
