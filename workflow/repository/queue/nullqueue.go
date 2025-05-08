package queue

import "time"

// NullInnerQueue null inner queue
type NullInnerQueue struct {
}

// Close Close
func (n *NullInnerQueue) Close() error {
	return nil
}

func (n *NullInnerQueue) CleanExecutionMessage(int) error {
	return nil
}

func (n *NullInnerQueue) SendInnerMessage(InnerMessage, *time.Time) error {
	return nil
}

func (n *NullInnerQueue) ReceiveInnerMessage() (<-chan InnerMessageBody, error) {
	return nil, nil
}

func (n *NullInnerQueue) SyncSchema() error {
	return nil
}
