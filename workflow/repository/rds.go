package execution

import "gorm.io/gorm"

type RDS interface {
	Session()
	SyncSchema() error
}

type Session = *gorm.Session
