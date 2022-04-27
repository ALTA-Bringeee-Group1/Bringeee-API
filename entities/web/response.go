package web

type SuccessListResponse struct {
	Status     string      `json:"status"`
	Code       int         `json:"code"`
	Error      interface{} `json:"error"`
	Links      interface{} `json:"links"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type SuccessResponse struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Error  interface{} `json:"error"`
	Links  interface{} `json:"links"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Error  string      `json:"error"`
	Links  interface{} `json:"links"`
}

type ValidationErrorResponse struct {
	Status string                `json:"status"`
	Code   int                   `json:"code"`
	Error  string                `json:"error"`
	Errors []ValidationErrorItem `json:"errors"`
	Links  interface{}           `json:"links"`
}

type ValidationErrorItem struct {
	Field string `json:"field"`
	Error string `json:"error"`
}
