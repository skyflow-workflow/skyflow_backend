package states

// Fail  失败节点
type Fail struct {
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

// JSONString fail data string sealize
func (fd FailData) String() string {
	jsonstr, _ := ToString(fd)
	return jsonstr
}

// GetFailData  return fail state data
func (s *Fail) GetFailData() FailData {
	fd := FailData{
		Error: s.Error,
		Cause: s.Cause,
	}
	return fd
}
