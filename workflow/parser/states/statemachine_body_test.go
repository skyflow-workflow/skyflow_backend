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
		definition  string
		bone        StateMachineBone
		expectError bool
	}{
		{
			name: "无State内容的状态机",
			definition: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"StartAt": "FirstState",
				"TimeoutSeconds": 10,
				"ExecutionInfoPath": "$.execution_info"
			}`,
			expectError: true,
		},
		{
			name:        "空状态机",
			definition:  "{}",
			expectError: true,
		},
		{
			name:        "无效的JSON格式",
			definition:  "{",
			expectError: true,
		},
		{
			name: "缺失必填字段",
			definition: `{
				"Comment": "An example of the Amazon States Language using a choice state.",
				"TimeoutSeconds": 10,
				"ExecutionInfoPath": "$.execution_info"
			}`,
			expectError: true,
		},
		{
			name: "测试无效Next的工作流流转内容",
			definition: `{
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
		{
			name: "正常工作流流转内容",
			definition: `{
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
			bone: StateMachineBone{
				StartAt: "FirstState",
				States: map[string]StateBone{
					"FirstState": {
						BaseBone: BaseBone{
							Type:    "Task",
							Name:    "FirstState",
							Next:    []string{"ChoiceState"},
							End:     false,
							Comment: "",
						},
					},
					"ChoiceState": {
						BaseBone: BaseBone{
							Type:    "Choice",
							Name:    "ChoiceState",
							Next:    []string{"FirstMatchState", "SecondMatchState", "DefaultState"},
							End:     false,
							Comment: "",
						},
					},
					"FirstMatchState": {
						BaseBone: BaseBone{
							Type:    "Succeed",
							Name:    "FirstMatchState",
							Next:    []string{},
							End:     true,
							Comment: "",
						},
					},
					"SecondMatchState": {
						BaseBone: BaseBone{
							Type:    "Succeed",
							Name:    "SecondMatchState",
							Next:    []string{},
							End:     true,
							Comment: "",
						},
					},
					"DefaultState": {
						BaseBone: BaseBone{
							Type:    "Fail",
							Name:    "DefaultState",
							Next:    []string{},
							End:     true,
							Comment: "",
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			smbody, err := NewStateMachineBodyFromString(tt.definition, StartDepth)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			bone := smbody.GetBone()
			assert.Equal(t, tt.bone, bone)
		})
	}
}
