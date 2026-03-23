package dto

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role    string `json:"role"`
}
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}