package executor

// common DB query fields for different scenarios, often used in execution and step queries.

// ExecutionFields commonly used fields for querying the Execution table.
var ExecutionFields = struct {
	// L1  least fields, only need to confirm if Execution exists,
	// {"id", "uuid", "status", "uri"}
	L1 []string
	// L2 basic fields, commonly used simple type fields
	// {"id", "uuid", "status", "max_execute_index", "execute_count", "header"}
	L2 []string
	// L3 extra fields compared to L2, includes JSON fields like definition/input/output
	// {"id", "uuid", "status", "max_execute_index", "execute_count", "header",
	// "definition"}
	L3 []string
	// L4 includes definition and input/output related fields for execution phase
	// {"id", "uuid", "status", "max_execute_index", "execute_count", "header",
	// "definition", "input", "output", "exception"},
	L4 []string
	//L5 查询所有字段
	L5 []string
}{
	L1: []string{"id", "uuid", "status", "uri"},
	L2: []string{"id", "uuid", "status", "max_execute_index", "execute_count", "header"},
	L3: []string{"id", "uuid", "status", "max_execute_index", "execute_count", "header",
		"definition"},
	L4: []string{"id", "uuid", "status", "max_execute_index", "execute_count", "header",
		"definition", "input", "output", "exception"},
	L5: []string{},
}

// StepFields commonly used fields for querying the Step structure.
var StepFields = struct {
	// L1 least fields, only need to confirm if Step exists
	// {"id", "execution_id", "status", "type", "name"}
	L1 []string
	// L2 basic simple type fields
	//{"id", "execution_id", "execute_index", "group_id", "group_index", "status", "name", "type", "execute_count"}
	L2 []string
	// L3 supports definition and syntax parsing
	//{"id", "execution_id", "execute_index", "group_id", "group_index", "status", "name", "type", "execute_count",
	// "resource", "definition"}
	L3 []string
	//L4 includes definition and input/output related fields for execution phase
	// {"id", "execution_id", "execute_index", "group_id", "status", "name", "type", "execute_count",
	// "resource", "definition", "input", "output", "data"}
	L4 []string
	// L5 query includes all fields, including context time and others
	L5 []string
}{
	L1: []string{"id", "execution_id", "status", "type", "name"},
	L2: []string{"id", "execution_id", "execute_index", "group_id", "group_index", "status", "name", "type", "execute_count"},
	L3: []string{"id", "execution_id", "execute_index", "group_id", "group_index", "status", "name", "type", "execute_count",
		"resource", "definition"},
	L4: []string{"id", "execution_id", "execute_index", "group_id", "group_index", "status", "name", "type", "execute_count",
		"resource", "definition", "input", "output", "data"},
	L5: []string{},
}

var StepFieldNames = struct {
	ID           string
	Name         string
	ExecutionID  string
	Type         string
	Status       string
	ExecuteIndex string
	ExecuteCount string
	GroupID      string
	GroupIndex   string
	Resource     string
	Definition   string
	Input        string
	Output       string
	Data         string
}{
	ID:           "id",
	Name:         "name",
	ExecutionID:  "execution_id",
	Type:         "type",
	Status:       "status",
	ExecuteIndex: "execute_index",
	ExecuteCount: "execute_count",
	GroupID:      "group_id",
	GroupIndex:   "group_index",
	Resource:     "resource",
	Definition:   "definition",
	Input:        "input",
	Output:       "output",
	Data:         "data",
}

var ExecutionFieldNames = struct {
	ID              string
	UUID            string
	FlowType        string
	Status          string
	MaxExecuteIndex string
	ExecuteCount    string
	Header          string
	Definition      string
	Input           string
	Output          string
	Exception       string
	URI             string
}{
	ID:              "id",
	UUID:            "uuid",
	FlowType:        "flow_type",
	Status:          "status",
	MaxExecuteIndex: "max_execute_index",
	ExecuteCount:    "execute_count",
	Header:          "header",
	Definition:      "definition",
	Input:           "input",
	Output:          "output",
	Exception:       "exception",
	URI:             "uri",
}
