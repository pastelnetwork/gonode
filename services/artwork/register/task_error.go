package register

type TaskError struct {
	error
}

func NewTaskError(err error) *TaskError {
	return &TaskError{
		error: err,
	}
}
