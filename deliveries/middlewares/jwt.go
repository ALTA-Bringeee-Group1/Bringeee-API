package middleware

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const secretJwt = "SECRET"

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte(secretJwt),
		SigningMethod: jwt.SigningMethodHS256.Name,
	})
}

func CreateToken(id int, name string, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	claims["name"] = name
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretJwt))
}

func ReadToken(token interface{}) (int, string, error) {
	tokenID := token.(*jwt.Token)
	claims := tokenID.Claims.(jwt.MapClaims)
	id := int(claims["id"].(float64))
	role := claims["role"].(string)
	return id, role, nil
}
