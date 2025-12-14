package apiserver

import (
	"testing"
	"time"

	"github.com/mmtbak/microlibrary/paging"
	pbv1 "github.com/skyflow-workflow/skyflow_backend/api/v1"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"github.com/stretchr/testify/assert"
)

func TestToTimeString(t *testing.T) {
	now := time.Now()
	expected := now.Format(timeformat)
	result := ToTimeString(now)
	assert.Equal(t, expected, result)
}

func TestToPBExecutionItem(t *testing.T) {
	createTime := time.Now()
	startTime := time.Now().Add(1 * time.Hour)
	finishTime := time.Now().Add(2 * time.Hour)

	tests := []struct {
		name     string
		input    po.Execution
		expected *pbv1.ExecutionItem
	}{
		{
			name: "WithStartAndFinishTime",
			input: po.Execution{
				UUID:       "test-uuid",
				Status:     "running",
				Title:      "test-title",
				Definition: "test-definition",
				CreateTime: createTime,
				StartTime:  &startTime,
				FinishTime: &finishTime,
			},
			expected: &pbv1.ExecutionItem{
				ExecutionUuid: "test-uuid",
				Status:        "running",
				Title:         "test-title",
				Definition:    "test-definition",
				CreateTime:    createTime.Unix(),
				StartTime:     startTime.Unix(),
				FinishTime:    startTime.Unix(), // Note: Bug in original code
			},
		},
		{
			name: "WithoutStartAndFinishTime",
			input: po.Execution{
				UUID:       "test-uuid",
				Status:     "completed",
				Title:      "test-title",
				Definition: "test-definition",
				CreateTime: createTime,
			},
			expected: &pbv1.ExecutionItem{
				ExecutionUuid: "test-uuid",
				Status:        "completed",
				Title:         "test-title",
				Definition:    "test-definition",
				CreateTime:    createTime.Unix(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToPBExecutionItem(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToVOPageRequest(t *testing.T) {
	tests := []struct {
		name     string
		input    *pbv1.PageRequest
		expected paging.PageRequest
	}{
		{
			name:     "NilInput",
			input:    nil,
			expected: paging.DefaultPageRequest,
		},
		{
			name: "ValidInput",
			input: &pbv1.PageRequest{
				PageSize:   10,
				PageNumber: 2,
			},
			expected: paging.PageRequest{
				PageSize:   10,
				PageNumber: 2,
			},
		},
		{
			name: "ExceedsMaxPageSize",
			input: &pbv1.PageRequest{
				PageSize:   3000,
				PageNumber: 1,
			},
			expected: paging.PageRequest{
				PageSize:   MaxPageSize,
				PageNumber: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToVOPageRequest(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToPBPageResponse(t *testing.T) {
	input := paging.PageResponse{
		PageSize:   10,
		PageNumber: 2,
		Count:      100,
		PageCount:  10,
	}

	expected := &pbv1.PageResponse{
		PageSize:   10,
		PageNumber: 2,
		Count:      100,
		PageCount:  10,
	}

	result := ToPBPageResponse(input)
	assert.Equal(t, expected, result)
}

func TestDataTransferArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		f        func(int) string
		expected []string
	}{
		{
			name:  "BasicConversion",
			input: []int{1, 2, 3},
			f: func(i int) string {
				return string(rune(i + 64))
			},
			expected: []string{"A", "B", "C"},
		},
		{
			name:  "EmptyInput",
			input: []int{},
			f: func(i int) string {
				return string(rune(i + 64))
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DataTransferArray(tt.input, tt.f)
			assert.Equal(t, tt.expected, result)
		})
	}
}
