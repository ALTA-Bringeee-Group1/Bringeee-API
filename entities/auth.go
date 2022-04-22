package entities

type AuthRequest struct {
	Email string `form:"email"`
	Password string `form:"password"`
}

type CustomerAuthResponse struct {
	Token string `json:"token"`
	User CustomerResponse
}

type DriverAuthResponse struct {
	Token string `json:"token"`
	User DriverResponse
}
type AdminAuthResponse struct {
	Token string `json:"token"`
	User AdminResponse
}