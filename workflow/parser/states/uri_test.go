package states

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestParseResourceURI(t *testing.T) {

	RegisterResource("aliyun")
	RegisterResource("aws")
	RegisterResource("arn")
	var testcases = []struct {
		resource string
		expected *ResourceURI
	}{
		{
			resource: "activity:unittest/add",
			expected: &ResourceURI{
				ResourceType: "activity",
				Resource:     "activity:unittest/add",
				Function:     "unittest/add",
			},
		},
		{
			resource: "activity:unittest/path/to/function/method",
			expected: &ResourceURI{
				ResourceType: "activity",
				Resource:     "activity:unittest/path/to/function/method",
				Function:     "unittest/path/to/function/method",
			},
		},
		{
			resource: "builtin:http/http_post",
			expected: &ResourceURI{
				Resource:     "builtin:http/http_post",
				ResourceType: "builtin",
				Function:     "http/http_post",
			},
		},
		{
			resource: "arn:FC:InvokeFunction",
			expected: &ResourceURI{
				Resource:     "arn:FC:InvokeFunction",
				ResourceType: "arn",
				Function:     "FC:InvokeFunction",
			},
		},
		{
			resource: "aws:arn:aws:states:::lambda:invoke",
			expected: &ResourceURI{
				Resource:     "aws:arn:aws:states:::lambda:invoke",
				ResourceType: "aws",
				Function:     "arn:aws:states:::lambda:invoke",
			},
		},
		{
			resource: "aliyun:acs:fc:::services/myService1.LATEST/functions/myFunction1",
			expected: &ResourceURI{
				Resource:     "aliyun:acs:fc:::services/myService1.LATEST/functions/myFunction1",
				ResourceType: "aliyun",
				Function:     "acs:fc:::services/myService1.LATEST/functions/myFunction1",
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.resource, func(t *testing.T) {
			resuri, err := ParseResource(tt.resource)
			assert.Equal(t, err, nil)
			assert.Equal(t, tt.expected, resuri)
		})
	}

}
