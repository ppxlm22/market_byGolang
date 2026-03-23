package repository

import "go_shopmarket/login/dto"

type Repository interface {
	GetUserByUsername(username string) (dto.UserResponse, string, error)
}
