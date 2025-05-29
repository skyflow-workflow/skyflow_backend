package queue

import (
	"fmt"
	"time"
)

// SimpleInnerQueue test inner queue
type SimpleInnerQueue struct {
	msgChan chan InnerMessage
}

// NewSimpleInnerQueue create a new simple inner queue
func NewSimpleInnerQueue() *SimpleInnerQueue {
	return &SimpleInnerQueue{
		msgChan: make(chan InnerMessage),
	}
}

// Close Close
func (t *SimpleInnerQueue) Close() error {
	if t.msgChan != nil {
		close(t.msgChan)
	}
	return nil
}

func (n *SimpleInnerQueue) CleanExecutionMessage(int) error {
	return nil
}

func (t *SimpleInnerQueue) SendInnerMessage(msg InnerMessageBody, sendTime *time.Time) error {
	newMsg := &TestInnerMessage{
		id:   fmt.Sprintf("%d", time.Now().UnixNano()),
		body: msg,
	}
	t.msgChan <- newMsg
	return nil
}

func (n *SimpleInnerQueue) ReceiveInnerMessage() (<-chan InnerMessage, error) {
	return n.msgChan, nil
}

func (n *SimpleInnerQueue) SyncSchema() error {
	return nil
}

type TestInnerMessage struct {
	id   string
	body InnerMessageBody
}

func (t *TestInnerMessage) ID() string {
	return fmt.Sprintf("%s", t.id)
}
func (t *TestInnerMessage) Body() InnerMessageBody {
	return t.body
}
func (t *TestInnerMessage) Ack() error {
	return nil
}
func (t *TestInnerMessage) Nack() error {
	return nil
}
