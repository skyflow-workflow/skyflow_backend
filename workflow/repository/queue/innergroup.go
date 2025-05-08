package queue

import (
	"fmt"
	"time"

	"github.com/mmtbak/microlibrary/config"
)

// InnerQueueGroup 内部消息队列组
// 包含正常队列与延时队列
type InnerQueueGroup struct {
	// Normal InnerMessage Queue  普通的内部消息队列
	_NormalQueue InnerMessageQueue
	// Deplay InnerMessage Queue  需要延时的内部消息队列
	_DelayQueue InnerMessageQueue
}

func NewInnerQueueGroupFromConfig(normal config.AccessPoint, delay config.AccessPoint) (*InnerQueueGroup, error) {

	var err error
	var normalqueue InnerMessageQueue
	var delayqueue InnerMessageQueue
	// queue define later
	if normal.Source == "" {
		err = fmt.Errorf("lack of config for inner queue")
		return nil, err
	}
	normalqueue, err = NewInnerMessageQueueFromConfig(normal)
	if err != nil {
		return nil, err
	}

	if delay.Source != "" {
		delayqueue, err = NewInnerMessageQueueFromConfig(delay)
		if err != nil {
			return nil, err
		}
	}

	groupqueue, err := NewInnerQueueGroup(normalqueue, delayqueue)
	if err != nil {
		return nil, err
	}
	return groupqueue, nil
}

// NewInnerQueueGroup create new inner queue group

func NewInnerQueueGroup(masterqueue InnerMessageQueue, delayqueue InnerMessageQueue) (*InnerQueueGroup, error) {

	if delayqueue != nil {
		if dbqueue, ok := delayqueue.(*DBMessageQueue); ok {
			dbqueue.StartPolling(masterqueue)
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
