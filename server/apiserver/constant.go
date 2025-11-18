package apiserver

import "github.com/skyflow-workflow/skyflow_backbend/workflow/vo"

var DefaultGetActivityRequest = vo.GetActivityTaskRequest{
	TimeoutSeconds: 10,
}
