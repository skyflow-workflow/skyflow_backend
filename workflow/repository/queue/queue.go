// package queue define the queue interface and message type for workflow
package queue

/*
 * @Author: mumangtao@gmail.com
 * @Date: 2020-11-05 18:28:22
 * @Last Modified by: mumangtao@gmail.com
 * @Last Modified time: 2020-11-05 18:34:09
 */

import (
	"encoding/json"
	"time"
)

// InnerMessageQueue  工作流内部执行的MQ
type InnerMessageQueue interface {
	// 发送消息
	SendInnerMessage(InnerMessageBody, *time.Time) error
	// 接收消息
	ReceiveInnerMessage() (<-chan InnerMessage, error)
	// 同步数据模型
	SyncSchema() error
	// 清理Execution相关的Message
	CleanExecutionMessage(int) error
	// Close 关闭MessageQueue
	Close() error
}

// TaskMessageQueue  工作流内部执行的MQ
type TaskMessageQueue interface {
	GetName() string
	// 同步数据模型
	SyncSchema() error
	// 发送消息
	SendTaskMessage(TaskMessageBody) error
	// 接收消息
	ReceiveTaskMessage() (<-chan TaskMessage, error)
	// Close 关闭MessageQueue
	Close() error
}

// InnerMessage 内部消息的队列接口
type InnerMessage interface {
	ID() string
	Body() InnerMessageBody
	Ack() error
	Nack() error
}

// TaskMessage 异步任务的队列接口
type TaskMessage interface {
	ID() string
	Body() TaskMessageBody
	Ack() error
	Nack() error
}

// InnerMessage 工作流内部的消息类型
type InnerMessageBody struct {
	ExecutionID int    //message execution id
	StepID      int    // message state id
	Class       string // class Execution、State
	Type        string // message type
	Data        string // data in message
}

func (imb *InnerMessageBody) Marshal() ([]byte, error) {
	return json.Marshal(imb)
}

func UnmarshalInnerMessageBody(data []byte) (*InnerMessageBody, error) {
	var imb InnerMessageBody
	err := json.Unmarshal(data, &imb)
	if err != nil {
		return nil, err
	}
	return &imb, nil
}

// AsyncTaskMessage 异步任务的消息类型
type TaskMessageBody struct {
	Type     string //Task Resource 类型，Resource 字段的前缀部分， qrn/arn 等等
	Resource string //Resource 的资源描述 "qrn:qcs:.xxxxxxxx:xxx"
	Input    string // Task 执行时的input
	Token    string // Task 执行时的token
	TaskData string // Task 执行时附带的数据， Execution创建时附带的
}

// MessageQueueErrorCode  errorcode
var MessageQueueErrorCode = struct {
	DBError         string
	DataError       string
	MessageNotFound string
}{
	DBError:         "DBError",
	DataError:       "DataError",
	MessageNotFound: "MessageNotFound",
}

// MessageClass Message class
var MessageClass = struct {
	Execution string
	State     string
	StepGroup string
}{
	Execution: "Execution",
	State:     "State",
	StepGroup: "StepGroup",
}
