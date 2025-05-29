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
func (mock MockInnerQueue) SendInnerMessage(msg InnerMessage, sendtime *time.Time) error {

	if sendtime == nil {
		sendtime = &time.Time{}
		*sendtime = time.Now()
	}
	s, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	msgstr := fmt.Sprintf("mock queue send message : send time : %s, message : %s\n", sendtime.String(), s)
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
