package mock

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestMockDB(t *testing.T) {

	client := GetMockDBClient()
	assert.NotEqual(t, client, nil)
}
