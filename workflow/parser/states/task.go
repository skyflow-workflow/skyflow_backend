package states

type Task struct {
	Name     string
	Resource string
}

func (t *Task) GetName() string {
	return t.Name
}

func (t *Task) Validate() string {
	return t.Resource
}
