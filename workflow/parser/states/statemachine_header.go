package states

// HeaderFieldNames ...
var HeaderFieldNames = struct {
	Version       string
	Type          string
	Comment       string
	QueryLanguage string
}{
	Version:       "Version",
	Type:          "Type",
	Comment:       "Comment",
	QueryLanguage: "QueryLanguage",
}

// StateMachineHeader ...
type StateMachineHeader struct {
	Version       string
	Type          string
	Comment       string
	QueryLanguage string
}

// Init ...
func (header *StateMachineHeader) Init() error {
	var err error
	if header.Version == "" {
		err = NewFieldPathError(ErrorInvalidFiledContent, HeaderFieldNames.Version)
		return err
	}
	if header.Type == "" {
		err = NewFieldPathError(ErrorInvalidFiledContent, HeaderFieldNames.Type)
		return err
	}
	// TODO  validate QueryLanguage later
	// if header.QueryLanguage == "" {
	// 	err = NewFieldError(ErrorInvalidFiledContent, HeaderFieldNames.QueryLanguage)
	// 	return err
	// }
	return nil
}
