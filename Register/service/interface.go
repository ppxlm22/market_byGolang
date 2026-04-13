package service
import (
	"go_shopmarket/register/dto"
	"go_shopmarket/register/repository"
)

type Service interface {
	RegisterUser(req dto.RegisterRequest) (*dto.RegisterResponse, error)
}

type service struct {
	repo repository.Repository
}