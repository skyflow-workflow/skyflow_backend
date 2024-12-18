package po

import "time"

// Namespace  资源namespace
type Namespace struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"not null; type:VARCHAR(255);unique; comment:name"` //namespace name
	Comment     string    `json:"comment" gorm:"type:MEDIUMTEXT"`
	GmtModified time.Time `json:"gmt_modified" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	GmtCreated  time.Time `json:"gmt_created" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}

// Function  每一个函数功能被称为一个 Function，用来描述一个特定的功能,有自己唯一的URI
type Function struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	NamespaceID int       `json:"namespace_id" gorm:"not null; uniqueIndex:uni_namespace_activity;INT(11)"`
	Name        string    `json:"name" gorm:"not null; uniqueIndex:uni_namespace_activity;type:VARCHAR(128)"`
	Type        string    `json:"type" gorm:"not null;type:VARCHAR(100)"`        // function type, activity/builtin/...
	URI         string    `json:"uri" gorm:"not null; unique;type:VARCHAR(256)"` // function uri
	Group       string    `json:"group" gorm:"not null;type:VARCHAR(128)"`       // function group
	Description string    `json:"description" gorm:"type:TEXT"`
	Parameters  string    `json:"parameters" gorm:"type:JSON"`              // paramters descritpion
	Status      string    `json:"status" gorm:"not null;type:VARCHAR(100)"` // active |disable
	GmtModified time.Time `json:"gmt_modified" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	GmtCreated  time.Time `json:"gmt_created" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}

// StateMachine 描述一个工作流的状态转化图
type StateMachine struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	NamespaceID int       `json:"namespace_id" gorm:"not null; uniqueIndex:uni_namespace_statemachine;INT(11)"`
	Type        string    `json:"type" gorm:"not null;type:VARCHAR(100)"`
	Name        string    `json:"name" gorm:"not nul; uniqueIndex:uni_namespace_statemachine;type:VARCHAR(128)"`
	URI         string    `json:"uri" gorm:"not null; unique;type:VARCHAR(256)"`
	Definition  string    `json:"definition" gorm:"not null;type:MEDIUMTEXT"`
	Comment     string    `json:"comment" gorm:"type:MEDIUMTEXT"`
	Status      string    `json:"status" gorm:"not null;type:VARCHAR(100)"`
	GmtModified time.Time `json:"gmt_modified" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	GmtCreated  time.Time `json:"gmt_created" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}
