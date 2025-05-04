package states

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func TestRunWait(t *testing.T) {

	var testcases = []struct {
		name        string
		wait        *Wait
		input       map[string]any
		waitTime    time.Time
		expectError error
	}{
		{
			name: "wait 5 seconds from Seconds",
			wait: &Wait{
				BaseState: &BaseState{
					Type: "Wait",
					Next: "LookupAddress2",
				},
				WaitBody: &WaitBody{
					Seconds: 5,
				},
			},
			input:       map[string]any{},
			waitTime:    time.Now().Add(5 * time.Second),
			expectError: nil,
		},
		{
			name: "wait 5 seconds from SecondsPath",
			wait: &Wait{
				BaseState: &BaseState{
					Type: "Wait",
					Next: "LookupAddress2",
				},
				WaitBody: &WaitBody{
					SecondsPath: "$.seconds",
				},
			},
			input: map[string]any{
				"test":    "hello",
				"test01":  "111",
				"seconds": 5,
			},
			waitTime:    time.Now().Add(5 * time.Second),
			expectError: nil,
		},
		{
			name: "wait 5.5 seconds from SecondsPath",
			wait: &Wait{
				BaseState: &BaseState{
					Type: "Wait",
					Next: "LookupAddress2",
				},
				WaitBody: &WaitBody{
					SecondsPath: "$.seconds",
				},
			},
			input: map[string]any{
				"test":    "hello",
				"test01":  "111",
				"seconds": 5.5,
			},
			waitTime:    time.Now().Add(5 * time.Second),
			expectError: fmt.Errorf("SecondPath value should be type int: [ $.seconds ]"),
		},
		{
			name: "wait 5 seconds from Timestamp",
			wait: &Wait{
				BaseState: &BaseState{
					Type: "Wait",
					Next: "LookupAddress2",
				},
				WaitBody: &WaitBody{
					Timestamp: "2020-12-31 19:56:35",
				},
			},
			input: map[string]any{
				"test":      "hello",
				"test01":    "111",
				"timestamp": "2020-12-31 19:56:35",
			},
			waitTime:    time.Date(2020, 12, 31, 19, 56, 35, 0, time.UTC),
			expectError: nil,
		},
		{
			name: "wait 5 seconds from TimestampPath",
			wait: &Wait{
				BaseState: &BaseState{
					Type: "Wait",
					Next: "LookupAddress2",
				},
				WaitBody: &WaitBody{
					TimestampPath: "$.timestamp",
				},
			},
			input: map[string]any{
				"test":      "hello",
				"test01":    "111",
				"timestamp": "2020-12-31 19:56:35",
			},
			waitTime:    time.Date(2020, 12, 31, 19, 56, 35, 0, time.UTC),
			expectError: nil,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			wt, err := tt.wait.GetWakeupTime(tt.input)
			if err != nil {
				t.Log(err.Error())
				assert.Equal(t, tt.expectError != nil, true)
				assert.Equal(t, err.Error(), tt.expectError.Error())
				return
			}
			assert.Equal(t, wt.Unix(), tt.waitTime.Unix())

		})
	}

}
