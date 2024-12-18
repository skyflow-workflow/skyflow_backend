package po

import "time"

// MessageQueue 消息队列db实现
type MessageQueue struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ExecutionID int       `json:"execution_id" gorm:"not null; index;type: INT(10)"` //message execution ID
	Topic       string    `json:"topic" gorm:"not null;type:VARCHAR(255)"`           // topic
	Type        string    `json:"type" gorm:"not null;type:VARCHAR(255)"`            // execution message type
	StepID      int       `json:"step_id" gorm:"not null; index;type:INT(10)"`       //  state id
	Data        string    `json:"data" gorm:"type:MEDIUMTEXT;DEFAULT:NULL"`          // message data
	Status      string    `json:"status" gorm:"not null;type:VARCHAR(100)"`          // message type
	Info        string    `json:"info" gorm:"not null; type:VARCHAR(255)"`           // message processinfomation
	SendTime    time.Time `json:"send_time" gorm:"not null;index;type:TIMESTAMP"`    // 消息发送时间，表示队列中发送消息的时间
	GmtModified time.Time `json:"gmt_modified" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	GmtCreated  time.Time `json:"gmt_created" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}
