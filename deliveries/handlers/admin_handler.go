package handlers

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/deliveries/helpers"
	middleware "bringeee-capstone/deliveries/middlewares"
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	userService "bringeee-capstone/services/user"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	userService *userService.UserService
}

func NewAdminHandler(service *userService.UserService) *AdminHandler {
	return &AdminHandler{
		userService: service,
	}
}

func (handler AdminHandler) DeleteDriver(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))
	token := c.Get("user")
	_, role, err := middleware.ReadToken(token)
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers/" + c.Param("id")}
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}

	// call delete service
	err = handler.userService.DeleteDriver(id)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data: map[string]interface{}{
			"id": id,
		},
	})
}

func (handler AdminHandler) UpdateDriverByAdmin(c echo.Context) error {

	// Bind request to user request
	userReq := entities.UpdateDriverByAdminRequest{}
	c.Bind(&userReq)
	id, _ := strconv.Atoi(c.Param("id"))

	// Get token
	token := c.Get("user")
	_, role, err := middleware.ReadToken(token)
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers/" + c.Param("id")}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	// avatar
	files := map[string]*multipart.FileHeader{}
	stnk_file, _ := c.FormFile("stnk_file")
	if stnk_file != nil {
		files["stnk_file"] = stnk_file
	}

	// Update via user service call
	userRes, err := handler.userService.UpdateDriverByAdmin(userReq, id, files)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   userRes,
	})
}

func (handler AdminHandler) GetAllDriver(c echo.Context) error {

	// Get token
	token := c.Get("user")
	_, role, err := middleware.ReadToken(token)

	// Translate query param to map of filters
	filters := []map[string]string{}
	name := c.QueryParam("name")
	if name != "" {
		filters = append(filters, map[string]string{
			"field":    "name",
			"operator": "=",
			"value":    name,
		})
	}
	gender := c.QueryParam("gender")
	if gender != "" {
		filters = append(filters, map[string]string{
			"field":    "gender",
			"operator": "=",
			"value":    gender,
		})
	}
	status := c.QueryParam("status")
	if status != "" {
		filters = append(filters, map[string]string{
			"field":    "status",
			"operator": "=",
			"value":    status,
		})
	}
	account_status := c.QueryParam("account_status")
	if account_status != "" {
		filters = append(filters, map[string]string{
			"field":    "account_status",
			"operator": "=",
			"value":    account_status,
		})
	}
	truck_type := c.QueryParam("truck_type")
	if truck_type != "" {
		filters = append(filters, map[string]string{
			"field":    "truck_type_id",
			"operator": "=",
			"value":    truck_type,
		})
	}

	// Sort parameter
	sorts := []map[string]interface{}{}
	sortName := c.QueryParam("sortName")
	sorts = append(sorts, map[string]interface{}{
		"field": "name",
		"desc":  map[string]bool{"1": true, "0": false}[sortName],
	})

	sortAge := c.QueryParam("sortAge")
	if sortAge != "" {
		switch sortAge {
		case "1":
			sorts = append(sorts, map[string]interface{}{
				"field": "age",
				"desc":  true,
			})
		case "0":
			sorts = append(sorts, map[string]interface{}{
				"field": "age",
				"desc":  false,
			})
		}
	}
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers?limit=" + c.QueryParam("limit") + "&page=" + c.QueryParam("page")}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	// pagination param
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, "Limit Parameter format is invalid", links))
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		links := map[string]string{"self": configs.Get().App.BaseURL}
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, "page Parameter format is invalid", links))
	}
	links["self"] = configs.Get().App.BaseURL + "/api/drivers?limit=" + c.QueryParam("limit") + "&page=" + c.QueryParam("page")

	// Get all drivers
	driversRes, err := handler.userService.FindAllDriver(limit, page, filters, sorts)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	// Get pagination data
	pagination, err := handler.userService.GetPaginationDriver(limit, page, filters)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}

	links["first"] = configs.Get().App.BaseURL + "/api/drivers?limit=" + c.QueryParam("limit") + "&page=1"
	links["last"] = configs.Get().App.BaseURL + "/api/drivers?limit=" + c.QueryParam("limit") + "&page=" + strconv.Itoa(pagination.TotalPages)
	if pagination.Page > 1 {
		links["prev"] = configs.Get().App.BaseURL + "/api/drivers?limit=" + c.QueryParam("limit") + "&page=" + strconv.Itoa(pagination.Page-1)
	}
	if pagination.Page < pagination.TotalPages {
		links["next"] = configs.Get().App.BaseURL + "/api/drivers?limit=" + c.QueryParam("limit") + "&page=" + strconv.Itoa(pagination.Page+1)
	}

	// success response
	return c.JSON(200, web.SuccessListResponse{
		Status:     "OK",
		Code:       200,
		Error:      nil,
		Links:      links,
		Data:       driversRes,
		Pagination: pagination,
	})
}

func (handler AdminHandler) GetSingleDriver(c echo.Context) error {
	// Get param and token
	id, err := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/drivers/" + c.Param("id")}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}
	token := c.Get("user")
	_, role, err := middleware.ReadToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}

	// Get eventdata
	driver, err := handler.userService.FindDriver(id)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	}
	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   driver,
	})
}

func (handler AdminHandler) GetSingleCustomer(c echo.Context) error {
	// Get param and token
	id, err := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": configs.Get().App.BaseURL + "/api/customers/" + c.Param("id")}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}
	token := c.Get("user")
	_, role, err := middleware.ReadToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}

	// Get eventdata
	user, err := handler.userService.FindCustomer(id)
	if err != nil {
		return helpers.WebErrorResponse(c, err, links)
	} else if user.ID == 0 {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Error:  "bad request",
			Links:  links,
		})
	}
	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   user,
	})
}
