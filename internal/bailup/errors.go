package bailup

type BailupError struct {
	message string
	inner   error
}

func NewBailupError(message string, inner error) *BailupError {
	return &BailupError{
		message: message,
		inner:   inner,
	}
}

func (e *BailupError) Error() string {
	return e.message
}

func (e *BailupError) Unwrap() error {
	return e.inner
}
