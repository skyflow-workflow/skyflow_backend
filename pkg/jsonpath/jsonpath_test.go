package jsonpath

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"trpc.group/trpc-go/trpc-go/log"
)

func TestJSONPathCompile(t *testing.T) {

	tests := []struct {
		name     string
		jsonpath string
		want     error
	}{
		{
			name:     "single level",
			jsonpath: "$.input",
			want:     nil,
		},
		{
			name:     "multi level",
			jsonpath: "$.input.input",
			want:     nil,
		},
		{
			name:     "start without $",
			jsonpath: "input.input.input",
			want:     nil,
		},
		{
			name:     "start with $",
			jsonpath: "$input.input.input",
			want:     fmt.Errorf("parse error at 2 in $input.input.input"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := JsonPathCompile(test.jsonpath)
			if test.want != nil {
				assert.Equal(t, test.want.Error(), err.Error())
			} else {
				assert.Equal(t, err, nil)
			}
		})
	}
}

func TestJSONPathGetValue(t *testing.T) {

	tests := []struct {
		name     string
		jsonpath string
		source   map[string]any
		value    any
		want     error
	}{
		{
			name:     "single level",
			jsonpath: "$.input",
			source:   map[string]any{"input": "inputvalue"},
			value:    "inputvalue",
			want:     nil,
		},
		{
			name:     "nil source",
			jsonpath: "$.input",
			source:   nil,
			value:    "inputvalue",
			want:     fmt.Errorf("key error: input not found in object"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := JsonPathGetValue(test.jsonpath, test.source)
			if err != nil {
				assert.Equal(t, test.want.Error(), err.Error())
			} else {
				assert.Equal(t, test.value, value)
			}
		})
	}
}

func TestJSONPathSetValue(t *testing.T) {

	tests := []struct {
		name     string
		jsonpath string
		source   any
		value    any
		result   any
		want     error
	}{
		{
			name:     "single level",
			jsonpath: "$.input",
			source:   map[string]any{"input": "inputvalue"},
			value:    "newvalue",
			result:   map[string]any{"input": "newvalue"},
			want:     nil,
		},
		{
			name:     "nil source",
			jsonpath: "$.input",
			source:   nil,
			value:    "newvalue",
			result:   nil,
			want:     fmt.Errorf("key error: input not found in object"),
		},
		{
			name:     "wrong jsonpath",
			jsonpath: "$.input.input",
			source:   map[string]any{"input": "inputvalue"},
			value:    "newvalue",
			result:   map[string]any{"input": map[string]any{"input": "newvalue"}},
			want:     fmt.Errorf("can not follow a string at '$.input'"),
		},
		{
			name:     "new jsonpath",
			jsonpath: "$.input2.input3",
			source:   map[string]any{"input": "inputvalue"},
			value:    "newvalue",
			result:   map[string]any{"input": "inputvalue", "input2": map[string]any{"input3": "newvalue"}},
			want:     nil,
		},
		{
			name:     "append jsonpath",
			jsonpath: "$.input.input",
			source:   map[string]any{"input": map[string]any{"input": "inputvalue"}},
			value:    "newvalue",
			result:   map[string]any{"input": map[string]any{"input": "newvalue"}},
			want:     nil,
		},
		{
			name:     "append2 jsonpath",
			jsonpath: "$.input2.input3",
			source:   map[string]any{"input": map[string]any{"input": "inputvalue"}},
			value:    "newvalue",
			result:   map[string]any{"input": map[string]any{"input": "inputvalue"}, "input2": map[string]any{"input3": "newvalue"}},
			want:     nil,
		},
		{
			name:     "append object",
			jsonpath: "$.input2.input3",
			source:   map[string]any{"input": map[string]any{"input": "inputvalue"}},
			value:    map[string]any{"input3": "input3value"},
			result:   map[string]any{"input": map[string]any{"input": "inputvalue"}, "input2": map[string]any{"input3": map[string]any{"input3": "input3value"}}},
			want:     nil,
		},
		{
			name:     "append array",
			jsonpath: "$.input2.input3",
			source:   map[string]any{"input": map[string]any{"input": "inputvalue"}},
			value:    []any{"input3value1", "input3value2"},
			result:   map[string]any{"input": map[string]any{"input": "inputvalue"}, "input2": map[string]any{"input3": []any{"input3value1", "input3value2"}}},
			want:     nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := JsonPathSetValue(test.jsonpath, test.source, test.value)
			if err != nil {
				if test.want != nil {
					log.Info(err.Error())
				}
				assert.Equal(t, test.want.Error(), err.Error())
			} else {
				assert.Equal(t, test.result, test.source)
			}
			t.Log(test.source)
		})
	}
}
