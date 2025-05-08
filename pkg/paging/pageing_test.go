package paging

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestPageing(t *testing.T) {

	req := PageRequest{
		PageNumber: 1,
		PageSize:   50,
	}
	limit, offset := req.Limit()
	assert.Equal(t, limit, 50)
	assert.Equal(t, offset, 0)
	resp := req.Response(103)
	assert.Equal(t, resp.PageCount, 3)
	assert.Equal(t, resp.PageNumber, 1)
	assert.Equal(t, resp.PageSize, 50)

	req = PageRequest{
		PageNumber: 2,
		PageSize:   50,
	}
	limit, offset = req.Limit()
	assert.Equal(t, limit, 50)
	assert.Equal(t, offset, 50)
	resp = req.Response(103)
	assert.Equal(t, resp.PageCount, 3)
	assert.Equal(t, resp.PageNumber, 2)
	assert.Equal(t, resp.PageSize, 50)

}
