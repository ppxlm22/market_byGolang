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
type RegisterRespone struct {
	Username string
        Email    string
        Message string
       
}
func (regDB *RegisterDB)ToModel() *RegisterRespone {
 return &RegisterRespone{
        Email: regDB.Email,
	Username: regDB.Username,
	Message: "success",
 } 
}
