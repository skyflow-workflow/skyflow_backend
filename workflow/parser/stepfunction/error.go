package stepfunction

import "errors"

// ErrorLackOfRequiredField ...
var (
	ErrorStateNameEmpty = errors.New("state name should not be empty")
)
