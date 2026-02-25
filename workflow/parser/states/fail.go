package states

// FailState  失败节点
type FailState struct {
	*BaseState `json:",inline" mapstructure:"BaseState"`
	*FailBody  `json:",inline" mapstructure:"FailBody"`
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

	data, err := StringToMap(definition)
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

// InitByMap Initialize Fail Content
func InitFailBodyByMap(body *FailBody, data map[string]interface{}) error {
	err := DecodeMapToStruct(data, body)
	if err != nil {
		return err
	}
	err = myValidate.Struct(body)
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

func (s *FailBody) GetDefinitionMap() (map[string]any, error) {
	data, err := DecodeStructToMap(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetDefinition returns the state definition for persistence.
// 使用 DecodeStructToMap 减少一次序列化，再展平为扁平 map、按拒绝名单删 key，保证输出 JSON 为扁平且可被 NewFailStateFromString 解析。
func (s *FailState) GetDefinition() (string, error) {
	baseData, err := s.BaseState.GetDefinitionMap()
	if err != nil {
		return "", err
	}
	bodyData, err := s.FailBody.GetDefinitionMap()
	if err != nil {
		return "", err
	}
	MapUpdate(baseData, bodyData)
	dataStr, err := ToString(baseData)
	return dataStr, err
}
