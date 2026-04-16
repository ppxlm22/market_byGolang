package dto

type RegisterRequest struct {
        Username string `json:"username" validate:"required"`
        Email    string `json:"email" validate:"required,email"`
        Password string `json:"password" validate:"required,min=6"`
}
type RegisterDB struct {
	Username string 
        Email    string
        Password string
}
type RegisterResponse struct {
	Username string
        Email    string
        Message string
       
}
func (regDB *RegisterDB)ToModel() *RegisterResponse {
        return &RegisterResponse{
                Email: regDB.Email,
                Username: regDB.Username,
                Message: "success",
        } 
}
