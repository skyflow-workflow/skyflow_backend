package states

import (
	"fmt"
	"strings"
)

var registedResourceType = make(map[string]struct{})

func init() {
	RegisterResource(ResourceType.Activity)
	RegisterResource(ResourceType.Builtin)
}

// RegisterResource ...
func RegisterResource(resource string) {
	registedResourceType[resource] = struct{}{}
}

// _ResourceType 资源类型

// ResourceSeparator ...
var ResourceSeparator = ":"

// ResourceType 资源类型
var ResourceType = struct {
	Builtin  string
	Activity string
}{
	Activity: "activity",
	Builtin:  "builtin",
}

// ResourceURI  resource uri defintion
type ResourceURI struct {
	Resource     string
	ResourceType string
	Function     string
}

// ParseResource parser path to ResourceURI
func ParseResource(resource string) (uri *ResourceURI, err error) {
	fields := strings.SplitN(strings.TrimSpace(resource), ResourceSeparator, 2)
	if len(fields) < 2 {
		err = fmt.Errorf("%w:invalid resource", ErrorInvalidFiledContent)
		return
	}
	restype := fields[0]
	if _, ok := registedResourceType[restype]; !ok {
		err = fmt.Errorf("%w:invalid resource type", ErrorInvalidFiledContent)
		return
	}
	resuri := &ResourceURI{
		Resource:     resource,
		ResourceType: fields[0],
		Function:     fields[1],
	}

	return resuri, nil
}
