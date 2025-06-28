package queue

import (
	"time"
)

// InnerQueueGroup inner queue group
type InnerQueueGroup struct {
	// Normal InnerMessage Queue
	_NormalQueue InnerMessageQueue
	// Deplay InnerMessage Queue
	_DelayQueue InnerMessageQueue
}

// NewInnerQueueGroup create new inner queue group

func NewInnerQueueGroup(masterqueue InnerMessageQueue, delayqueue InnerMessageQueue) (*InnerQueueGroup, error) {

	if delayqueue != nil {
		if dbqueue, ok := delayqueue.(*DBMessageQueue); ok {
			dbqueue.SetForwardQueue(masterqueue)
			dbqueue.StartPolling()
		}
	}

	qg := &InnerQueueGroup{
		_NormalQueue: masterqueue,
		_DelayQueue:  delayqueue,
	}
	return qg, nil
}

// SendInnerMessage  Select suitable queue send message
func (q *InnerQueueGroup) SendInnerMessage(msg InnerMessageBody, sendTime *time.Time) error {
	var err error
	now := time.Now()
	if sendTime != nil && now.Before(*sendTime) && q._DelayQueue != nil {
		err = q._DelayQueue.SendInnerMessage(msg, sendTime)
	} else {
		err = q._NormalQueue.SendInnerMessage(msg, sendTime)
	}
	return err
}

// ReceiveInnerMessage  receive innermessage from group queue
func (q *InnerQueueGroup) ReceiveInnerMessage() (<-chan InnerMessage, error) {
	var err error
	nchan, err := q._NormalQueue.ReceiveInnerMessage()
	return nchan, err
}

// CleanExecution 清理一个execution的 message
func (q *InnerQueueGroup) CleanExecutionMessage(execution_id int) error {

	var err error
	err = q._NormalQueue.CleanExecutionMessage(execution_id)
	if err != nil {
		return err
	}
	if q._DelayQueue != nil {
		err = q._DelayQueue.CleanExecutionMessage(execution_id)
		if err != nil {
			return err
		}
	}
	return nil
}

// Close  close group queue
func (q *InnerQueueGroup) Close() error {
	var err error
	err = q._NormalQueue.Close()
	if err != nil {
		return err
	}
	if q._DelayQueue != nil {
		err = q._DelayQueue.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// SyncSchema
func (q *InnerQueueGroup) SyncSchema() error {

	var err error
	err = q._NormalQueue.SyncSchema()
	if err != nil {
		return err
	}
	if q._DelayQueue != nil {
		err = q._DelayQueue.SyncSchema()
		if err != nil {
			return err
		}
	}
	return nil
}
