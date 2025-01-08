package states

type BaseState struct {
	Name string `json:"Name,omitempty"`
	Type string `json:"Type,omitempty"`
	End  bool   `json:"End,omitempty"`
	Next string `json:"Next,omitempty"`
}

func (s *BaseState) GetName() string {
	return s.Name
}

func (s *BaseState) SetName(name string) {
	s.Name = name
}

func (s *BaseState) GetType() string {
	return s.Type
}

func (s *BaseState) Init() error {
	return nil
}
