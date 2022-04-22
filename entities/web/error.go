package web

type WebError struct {
	Code    int
	Message string
}

func (err WebError) Error() string {
	return err.Message
}

type ValidationError struct {
	Code    int
	Message string
	Errors  []ValidationErrorItem
}

func (err ValidationError) Error() string {
	return err.Message
}
