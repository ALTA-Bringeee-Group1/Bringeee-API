package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/v4"
)

func main() {
	config := config.Get()
	db := utilities.NewMysqlGorm(config)
	utilities.Migrate(db)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
	}))
}
