package states

// FailState  失败节点
type FailState struct {
	*BaseState
	*FailBody
}

type FailBody struct {
	Abort bool   `mapstructure:"Abort"`
	Cause string `mapstructure:"Cause"`
	Error string `mapstructure:"Error"`
}

// FailData  Fail State Data
type FailData struct {
	Cause string `json:"cause"`
	Error string `json:"error"`
}

// JSONString fail data string serialization
func (fd FailData) String() string {
	jsonstr, _ := ToString(fd)
	return jsonstr
}

// NewFailStateFromString NewFailStateFromString
func NewFailStateFromString(definition string) (state *FailState, err error) {

	data, err := ToMap(definition)
	if err != nil {
		return nil, err
	}
	state, err = NewFailStateFromMap(data)
	return state, err
}

// NewFailStateFromMap Create New Fail State
func NewFailStateFromMap(data map[string]interface{}) (state *FailState, err error) {
	// basestate
	bs, err := NewBaseStateFromMap(data)
	if err != nil {
		return state, err
	}
	bs.End = true

	body := &FailBody{}
	err = InitFailBodyByMap(body, data)
	if err != nil {
		return state, err
	}
	state = &FailState{
		BaseState: bs,
		FailBody:  body,
	}
	return
}

func NewFailBodyFromMap(data map[string]interface{}) (*FailBody, error) {
	var err error
	failbody := &FailBody{}
	err = InitFailBodyByMap(failbody, data)
	if err != nil {
		return nil, err
	}
	return failbody, nil
}

// InitByMap Inititalize Fail Content
func InitFailBodyByMap(body *FailBody, data map[string]interface{}) error {
	err := MapStructDecode(data, body)
	if err != nil {
		return err
	}
	err = myvalidate.Struct(body)
	if err != nil {
		return err
	}
	return nil
}

// GetFailData  return fail state data
func (s *FailState) GetFailData() FailData {
	fd := FailData{
		Error: s.Error,
		Cause: s.Cause,
	}
	return fd
}

func (s *FailState) GetBaseState() *BaseState {
	return s.BaseState
}
