package po

import "time"

// Namespace  资源namespace
type Namespace struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string    `json:"name" gorm:"not null; type:VARCHAR(255);unique; comment:name"` //namespace name
	Comment    string    `json:"comment" gorm:"type:MEDIUMTEXT"`
	UpdateTime time.Time `json:"update_time" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	CreateTime time.Time `json:"create_time" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}

// Activity  每一个函数功能被称为一个 Activity,有自己唯一的URI
type Activity struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	NamespaceID int       `json:"namespace_id" gorm:"not null; uniqueIndex:uni_namespace_activity;INT(11)"`
	Name        string    `json:"name" gorm:"not null; uniqueIndex:uni_namespace_activity;type:VARCHAR(128)"`
	Type        string    `json:"type" gorm:"not null;type:VARCHAR(100)"`        // function type, activity/builtin/...
	URI         string    `json:"uri" gorm:"not null; unique;type:VARCHAR(256)"` // function uri
	Comment     string    `json:"comment" gorm:"type:TEXT"`
	Parameters  string    `json:"parameters" gorm:"type:TEXT"`              // paramters descritpion
	Status      string    `json:"status" gorm:"not null;type:VARCHAR(100)"` // active |disable
	UpdateTime  time.Time `json:"update_time" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	CreateTime  time.Time `json:"create_time" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}

// StateMachine 描述一个工作流的状态转化图
type StateMachine struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	NamespaceID int       `json:"namespace_id" gorm:"not null; uniqueIndex:uni_namespace_statemachine;INT(11)"`
	Name        string    `json:"name" gorm:"not nul; uniqueIndex:uni_namespace_statemachine;type:VARCHAR(128)"`
	URI         string    `json:"uri" gorm:"not null; unique;type:VARCHAR(256)"`
	Definition  string    `json:"definition" gorm:"not null;type:MEDIUMTEXT"`
	Comment     string    `json:"comment" gorm:"type:MEDIUMTEXT"`
	Status      string    `json:"status" gorm:"not null;type:VARCHAR(100)"`
	UpdateTime  time.Time `json:"update_time" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	CreateTime  time.Time `json:"create_time" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}
