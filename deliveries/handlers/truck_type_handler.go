package handlers

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/entities/web"
	truckTypeService "bringeee-capstone/services/truck_type"
	"fmt"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type TruckTypeHandler struct {
	truckTypeService truckTypeService.TruckTypeServiceInterface
}

func NewTruckTypeHandler(service truckTypeService.TruckTypeServiceInterface) *TruckTypeHandler {
	return &TruckTypeHandler{
		truckTypeService: service,
	}
}


/*
 * Handler: Find all Truck Type
 * -------------------------------
 * Mengambil data truckType berdasarkan filters dan sorts berbentuk query params
 * Public: GET /api/truck_type/
 */
func (handler TruckTypeHandler) Index(c echo.Context) error {

	// Url path parameter & link hateoas
	links := map[string]string {}
	
	links["self"] = fmt.Sprintf("%s/api/truck_type", configs.Get().App.ENV)

	// service call
	truckTypeRes, err := handler.truckTypeService.FindAll(0, 0, []map[string]string{}, []map[string]interface{}{})
	if err != nil {
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
			Code: 500,
			Error: "Server error: " + err.Error(),
			Links: links,
		})
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusOK,
		Error: nil,
		Links: links,
		Data: truckTypeRes,
	})
}