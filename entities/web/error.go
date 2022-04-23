package web

import "bringeee-capstone/configs"

type WebError struct {
	Code    int
	ProductionMessage string
	DevelopmentMessage string
	Message string
}

func (err WebError) Error() string {
	return map[bool]string{ true: err.ProductionMessage, false: err.DevelopmentMessage }[configs.Get().App.ENV == "production"]
}

type ValidationError struct {
	Code    int
	Message string
	Errors  []ValidationErrorItem
}

func (err ValidationError) Error() string {
	return err.Message
}
