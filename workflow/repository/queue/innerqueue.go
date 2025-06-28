package queue

import (
	"github.com/mmtbak/microlibrary/mq"
)

// NewInnerMessageQueueFromConfig create new inner message queue
func NewInnerMessageQueueFromConfig(dsn string) (InnerMessageQueue, error) {

	var err error
	basequeue, err := mq.NewMessageQueue(dsn)
	if err != nil {
		return nil, err
	}
	innerMQ := NewMQInnerMessageQueue(basequeue)
	return innerMQ, nil

}
