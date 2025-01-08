package states

type WaitBody struct {
	Seconds int `json:"Seconds"`
}

type Wait struct {
	*BaseState
	WaitBody
}
