package cmd

// FrameWorkAPIServer defines the interface for starting and stopping the framework API server.
// for diffrent framework, the implementation is different
type FrameWorkAPIServer interface {
	Start() error
	Stop() error
}
