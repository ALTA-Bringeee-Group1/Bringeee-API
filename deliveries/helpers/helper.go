package helpers

import (
	"bringeee-capstone/entities/web"

	"golang.org/x/crypto/bcrypt"
)

func ResponseFailed(message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  "failed",
		"message": message,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func MakeErrorResponse(status string, code int, err string, links map[string]string) web.ErrorResponse {
	return web.ErrorResponse{
		Status: status,
		Code:   code,
		Error:  err,
		Links:  links,
	}
}
