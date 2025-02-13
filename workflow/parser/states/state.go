package states

// State ...
type State interface {
	Init() error
	GetName() string
	SetName(name string)
	GetType() string
	GetBone() StateBone
}
