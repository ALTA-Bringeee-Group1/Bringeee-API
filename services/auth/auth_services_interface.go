package auth

import (
	"bringeee-capstone/entities"
)

type AuthServiceInterface interface {
	Login(AuthReq entities.AuthRequest) (interface{}, error)
	Me(ID int, token interface{}) (interface{}, error)
}
