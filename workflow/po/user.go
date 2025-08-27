package po

import "time"

// other persistent objects for workflow

// UserStepData user stored data for a state
type UserStepData struct {
	ID     int `json:"id" gorm:"primaryKey;autoIncrement"`
	StepID int `json:"step_id" gorm:"not null; unique; type:INT(11)"`
	// user store data for a step, string format
	StoreData string `json:"store_data" gorm:"type:MEDIUMTEXT"`
	// user store reference for a step, json format
	Reference  string    `json:"reference" gorm:"type:JSON"`
	ModifyTime time.Time `json:"gmt_modified" gorm:"<-:create update;autoUpdateTime;type:TIMESTAMP" `
	CreateTime time.Time `json:"gmt_created" gorm:"<-:create;autoCreateTime;type:TIMESTAMP"`
}
