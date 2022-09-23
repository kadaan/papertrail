package errors

type CommandError struct {
	message string
}

var _ error = &CommandError{}

func NewCommandError(message string) *CommandError {
	return &CommandError{
		message: message,
	}
}

func (e *CommandError) Error() string {
	return e.message
}
