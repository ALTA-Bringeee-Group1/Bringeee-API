package entities

type AuthRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type CustomerAuthResponse struct {
	Token string           `json:"token"`
	User  CustomerResponse `json:"user"`
}

type DriverAuthResponse struct {
	Token string         `json:"token"`
	User  DriverResponse `json:"user"`
}
type AdminAuthResponse struct {
	Token string        `json:"token"`
	User  AdminResponse `json:"user"`
}
