package states

type ChoiceBody struct {
	Choices []Choice `json:"Choices"`
}

type Choice struct {
	*BaseState
	ChoiceBody
}
