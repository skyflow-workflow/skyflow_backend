package decoder

type FieldError struct {
	Field string
	Error string
	Path  string
}
