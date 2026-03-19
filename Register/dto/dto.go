package dto

type RegisterRequest struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
}
type registerDB struct {
	Username string
        Email    string
        Password string
}
