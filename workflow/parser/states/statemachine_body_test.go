package states

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseStateMachineBody 测试解析状态机体的功能
// 包括正常情况和边界情况
func TestParseStateMachineBody(t *testing.T) {
	var testcases = []struct {
		name        string
		body        string
		expectError bool
	}{
		{
			name: "正常解析状态机体的JSON",
			body: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": 10,
				"ExecutionInfoPath": "$.execution_info"
			}`,
			expectError: true,
		},
		{
			name:        "空JSON",
			body:        "{}",
			expectError: true,
		},
		{
			name:        "无效的JSON格式",
			body:        "{",
			expectError: true,
		},
		{
			name: "缺失必填字段",
			body: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"TimeoutSeconds": 10,
				"ExecutionInfoPath": "$.execution_info"
			}`,
			expectError: true,
		},
		{
			name: "测试工作流流转内容",
			body: `{
				"Comment": "Test workflow transitions",
				"StartAt": "FirstState",
				"States": {
					"FirstState": {
						"Type": "Task",
						"Resource": "arn:aws:states:us-east-1:123456789012:activity:HelloWorld",
						"Next": "ChoiceState"
					},
					"ChoiceState": {
						"Type": "Choice",
						"Choices": [
							{
								"Variable": "$.foo",
								"NumericEquals": 1,
								"Next": "FirstMatchState"
							},
							{
								"Variable": "$.foo",
								"NumericEquals": 2,
								"Next": "SecondMatchState"
							}
						],
						"Default": "DefaultState"
					},
					"FirstMatchState": {
						"Type": "Succeed"
					},
					"SecondMatchState": {
						"Type": "Succeed"
					},
					"DefaultState": {
						"Type": "Fail",
						"Error": "DefaultStateError",
						"Cause": "No Matches!"
					}
				}
			}`,
			expectError: false,
		},
		{
			name: "测试无效的工作流流转内容",
			body: `{
				"Comment": "Invalid workflow transitions",
				"StartAt": "FirstState",
				"States": {
					"FirstState": {
						"Type": "Task",
						"Resource": "arn:aws:states:us-east-1:123456789012:activity:HelloWorld",
						"Next": "NonExistentState"
					}
				}
			}`,
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewStateMachineBodyFromString(tc.body, StartDepth)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
