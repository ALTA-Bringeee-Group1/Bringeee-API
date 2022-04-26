package auth

import (
	"bringeee-capstone/entities"
)

type AuthServiceInterface interface {
	Login(AuthReq entities.AuthRequest) (interface{}, error)
	Me(Id int, token interface{}) (interface{}, error)
}
