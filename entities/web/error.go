package web

import "bringeee-capstone/configs"

type WebError struct {
	Code    int
	ProductionMessage string
	DevelopmentMessage string
	Message string
}

func (err WebError) Error() string {
	if err.ProductionMessage == "" && err.DevelopmentMessage == "" {
		return err.Message
	} 
	return map[bool]string{ true: err.ProductionMessage, false: err.DevelopmentMessage }[configs.Get().App.ENV == "production"]
}

type ValidationError struct {
	Code    			int
	ProductionMessage 	string
	DevelopmentMessage 	string
	Errors  []ValidationErrorItem
}

func (err ValidationError) Error() string {
	return map[bool]string{ true: err.ProductionMessage, false: err.DevelopmentMessage }[configs.Get().App.ENV == "production"]
}
