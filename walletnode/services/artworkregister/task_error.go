package artworkregister

// TaskError represents an error of the task.
type TaskError struct {
	error
}

// NewTaskError returns a new TaskError instance.
func NewTaskError(err error) *TaskError {
	return &TaskError{
		error: err,
	}
}
