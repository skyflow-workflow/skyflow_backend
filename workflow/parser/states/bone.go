package states

// Bone is the framework structure of the workflow, common used for draw the workflow graph.
// It is also used for api to understand the workflow. and fetch the workflow execution status.

// BaseBone   common state type bone
type BaseBone struct {
	Type    string   `json:"Type"`
	Name    string   `json:"Name"`
	Next    []string `json:"Next"`
	End     bool     `json:"End"`
	Comment string   `json:"Comment"`
}

// StateBone for map /parallel state bone
type StateBone struct {
	BaseBone
	*StateMachineBone
	Branches []StateBone `json:"Branches"`
}

// StateMachineBone  StateMachineBone
type StateMachineBone struct {
	StartAt string               `json:"StartAt"`
	States  map[string]StateBone `json:"States"`
}
