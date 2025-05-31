package queue

import (
	"fmt"
	"log/slog"
	"sync"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestSimpleQueue(t *testing.T) {
	queue := NewSimpleInnerQueue()

	var msgs = []InnerMessageBody{
		{1, 1, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{1, 2, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{2, 2, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{2, 3, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{2, 4, "Step", "StepInit", `{"testkey":"testvalue"}`},
		{2, 5, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
	}

	var wg = &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		msgChan, err := queue.ReceiveInnerMessage()
		assert.Equal(t, err, nil)
		msgCount := 0
		for {
			msg, ok := <-msgChan
			if !ok {
				break
			}
			msgCount++
			slog.Info("receive message", "msg", msg)
		}
		slog.Info(fmt.Sprintf("receive message count: %d", msgCount))
		wg.Done()
		assert.Equal(t, msgCount, len(msgs))
	}()

	var err error
	for _, msg := range msgs {
		err = queue.SendInnerMessage(msg, nil)
		assert.Equal(t, err, nil)
		slog.Info("send message", "msg", msg)
	}
	slog.Info(fmt.Sprintf(
		"all messages sent, send message count: %d", len(msgs)),
	)
	err = queue.Close()
	assert.Equal(t, err, nil)
	wg.Wait()

}
