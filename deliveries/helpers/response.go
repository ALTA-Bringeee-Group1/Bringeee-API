package helpers

import (
	"bringeee-capstone/entities/web"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

func WebErrorResponse(c echo.Context, err error, links map[string]string) error {
	if reflect.TypeOf(err).String() == "web.WebError" {
		webErr := err.(web.WebError)
		return c.JSON(webErr.Code, web.ErrorResponse{
			Status: "ERROR",
			Code: webErr.Code,
			Error: webErr.Error(),
			Links: links,
		})
	}
	return c.JSON(http.StatusInternalServerError, web.ErrorResponse{
		Status: "ERROR",
		Code: http.StatusInternalServerError,
		Error: "Server Error",
		Links: links,
	})

}