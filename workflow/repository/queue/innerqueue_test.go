package queue

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestCreateInnerQueue(t *testing.T) {

	dsn := "kafka://localhost:9092/innerqueue_test?consumergroup=innerqueue_test_group"
	_, err := NewInnerMessageQueueFromConfig(dsn)
	assert.Equal(t, nil, err)
}
