package repository
import (
	"go_shopmarket/register/dto"
)

type Repository interface {
	CheckUserExists(username, email string) (bool, error)
	Register(req dto.RegisterRequest) (*dto.RegisterDB, error)

}