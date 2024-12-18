// persistance object for workflow
package po

// GetExecutionTables get tables for execution
func GetExecutionTables() []interface{} {

	return []interface{}{
		new(Execution),
		new(ExecutionShade),
		new(State),
		new(ActivityTask),
		new(ExecutionEvent),
		new(TaskToken),
		new(StateGroup),
	}
}

// GetEventTables get tables for exporter event
func GetEventTables() []interface{} {
	return []interface{}{
		new(ExecutionEvent),
	}
}

// GetTemplateTables get tables for workflow template
func GetTemplateTables() []interface{} {
	return []interface{}{
		new(Namespace),
		new(Function),
		new(StateMachine),
	}
}

// GetMQTables get tables for message queue
func GetMQTables() []interface{} {
	return []interface{}{
		new(MessageQueue),
	}
}
