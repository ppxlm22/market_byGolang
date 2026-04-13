package dto

type RegisterRequest struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
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
