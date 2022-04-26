package handlers

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/entities/web"
	regionService "bringeee-capstone/repositories/region"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type RegionHandler struct {
	regionService regionService.RegionRepositoryInterface
}

func NewRegionHandler(service regionService.RegionRepositoryInterface) *RegionHandler {
	return &RegionHandler{
		regionService: service,
	}
}

/*
 * List Province
 * ----------------------------------
 * Endpoint daftar provinsi 
 * Public GET /api/provinces
 */
func (handler RegionHandler) IndexProvince (c echo.Context) error {

	links := map[string]string{}
	links["self"] = fmt.Sprintf("%s/api/provinces/", strings.TrimSuffix(configs.Get().App.BaseURL, "/"))

	provinces, err := handler.regionService.FindAllProvince([]map[string]interface{}{})
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
			Code: http.StatusInternalServerError,
			Error: "server error",
			Links: links,
		})
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusOK,
		Data: provinces,
		Error: nil,
		Links: links,
	})
}

/*
 * List City
 * ----------------------------------
 * Endpoint daftar kota 
 * Public GET /api/provinces/{provinceID}/cities
 */
func (handler RegionHandler) IndexCity (c echo.Context) error {

	links := map[string]string{}
	links["self"] = fmt.Sprintf("%s/api/provinces/%s/cities/", 
		strings.TrimSuffix(configs.Get().App.BaseURL, "/"),
		c.Param("provinceID"),
	)

	provinceID, err := strconv.Atoi(c.Param("provinceID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code: 400,
			Error: "Invalid parameter",
			Links: links,
		})
	}


	cities, err := handler.regionService.FindAllCity(provinceID, []map[string]interface{}{})
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
			Code: http.StatusInternalServerError,
			Error: "server error",
			Links: links,
		})
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusOK,
		Data: cities,
		Error: nil,
		Links: links,
	})
}

/*
 * List City
 * ----------------------------------
 * Endpoint daftar kecamatan
 * Public GET /api/province/{provinceID}/cities/{cityID}/districts
 */
func (handler RegionHandler) IndexDistrict (c echo.Context) error {

	links := map[string]string{}
	links["self"] = fmt.Sprintf("%s/api/provinces/%s/cities/%s/districts", 
		strings.TrimSuffix(configs.Get().App.BaseURL, "/"),
		c.Param("provinceID"),
		c.Param("cityID"),
	)

	cityID, err := strconv.Atoi(c.Param("cityID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Status: "ERROR",
			Code: 400,
			Error: "Invalid parameter",
			Links: links,
		})
	}

	cities, err := handler.regionService.FindAllDistrict(cityID, []map[string]interface{}{})
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
			Code: http.StatusInternalServerError,
			Error: "server error",
			Links: links,
		})
	}

	return c.JSON(http.StatusOK, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusOK,
		Data: cities,
		Error: nil,
		Links: links,
	})
}