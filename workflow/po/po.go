// persistent object for workflow
package po

// GetExecutionTables get tables for execution
func GetExecutionTables() []any {

	return []any{
		new(Execution),
		new(Step),
		new(ActivityTask),
		new(ExecutionEvent),
		new(TaskToken),
		new(StepGroup),
	}
}

// GetEventTables get tables for exporter event
func GetEventTables() []any {
	return []any{
		new(ExecutionEvent),
	}
}

// GetTemplateTables get tables for workflow template
func GetTemplateTables() []any {
	return []any{
		new(Namespace),
		new(Activity),
		new(StateMachine),
	}
}

// GetMQTables get tables for message queue
func GetMQTables() []any {
	return []any{
		new(MessageQueue),
	}
}
