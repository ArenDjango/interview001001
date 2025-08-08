package handler

type NothingFound struct {
	originalError error
}

func (e *NothingFound) Error() string {
	return "Nothing found: " + e.originalError.Error()
}

func NewNothingFound(err error) error {
	return &NothingFound{originalError: err}
}
