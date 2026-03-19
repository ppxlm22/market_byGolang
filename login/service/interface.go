package service
import (
	"go_shopmarket/login/repository"
	"go_shopmarket/login/dto"
)

type Service interface {
	LoginUser(req dto.LoginRequest) (string, error)
}
type service struct {
	repo repository.Repository
}
