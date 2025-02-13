package states

// ChoiceBody ...
type ChoiceBody struct {
	Choices []Choice `json:"Choices"`
}

// Choice ...
type Choice struct {
	*BaseState
	ChoiceBody
}
