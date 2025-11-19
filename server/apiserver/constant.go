package apiserver

import "github.com/skyflow-workflow/skyflow_backend/workflow/vo"

var DefaultGetActivityRequest = vo.GetActivityTaskRequest{
	TimeoutSeconds: 10,
}
