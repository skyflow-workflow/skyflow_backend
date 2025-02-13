package states

// WaitBody ...
type WaitBody struct {
	Seconds int `json:"Seconds"`
}

// Wait ...
type Wait struct {
	*BaseState
	WaitBody
}
