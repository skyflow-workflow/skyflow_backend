package queue

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

// MockInnerQueue mock innerqueue
type MockInnerQueue struct{}

func (mock MockInnerQueue) GetName() string {
	return "mockqueue"
}

// SendInnerMessage SendInnerMessage
func (mock MockInnerQueue) SendInnerMessage(msg InnerMessage, desttime time.Time) error {

	s, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	msgstr := fmt.Sprintf("mock queue send message : send time : %s, message : %s\n", desttime.String(), s)
	slog.Info(msgstr)
	return nil
}

// ReceiveInnerMessage ReceiveInnerMessage
func (mock MockInnerQueue) ReceiveInnerMessage() (<-chan InnerMessageBody, error) {
	return nil, nil
}

// Close Close
func (mock MockInnerQueue) Close() error {
	return nil
}

// CleanExecution 清理一个execution的 message
// 在这里， 这个是一个空方法
func (mock MockInnerQueue) CleanExecutionMessage(int) error {
	return nil
}

func (mock MockInnerQueue) SyncSchema() error {
	return nil
}

// MockTaskQueueGroup mock task queue group
type MockTaskQueueGroup struct{}

// CreateQueue CreateQueue
func (v MockTaskQueueGroup) CreateQueue(string) (TaskMessageQueue, error) {
	return nil, nil
}

// SendTaskMessage SendTaskMessage
func (v MockTaskQueueGroup) SendTaskMessage(name string, msg TaskMessageBody) error {
	s, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	msgstr := fmt.Sprintf("mock queue send message, send time : %s, message : %s\n", time.Now().String(), s)
	slog.Info(msgstr)
	return nil
}

// ReceiveTaskMessage ReceiveTaskMessage
func (v MockTaskQueueGroup) ReceiveTaskMessage(string) (<-chan TaskMessage, error) {
	return nil, nil
}

// Close Close
func (v MockTaskQueueGroup) Close() error {
	return nil
}

func (v MockTaskQueueGroup) SyncSchema() error {
	return nil
}
