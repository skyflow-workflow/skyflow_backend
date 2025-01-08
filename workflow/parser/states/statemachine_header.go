package states

var (
	DefaultStateMachineHeader = StateMachineHeader{
		Version: "1.0",
		Type:    "stepfunction",
	}
)

type StateMachineHeader struct {
	Version string
	Type    string
	Comment string
}

func (header *StateMachineHeader) Init() error {
	if header.Version == "" {
		header.Version = DefaultStateMachineHeader.Version
	}
	if header.Type == "" {
		header.Type = DefaultStateMachineHeader.Type
	}
	return nil
}
