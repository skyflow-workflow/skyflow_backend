package states

type TaskBody struct {
	Resource       string `json:"Resource"`
	TimeoutSeconds int    `json:"TimeoutSeconds"`
}

type Task struct {
	*BaseState
	TaskBody
}

func (t *Task) Validate() string {
	return t.Resource
}
