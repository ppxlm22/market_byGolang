package register

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func Register_Service(req RegisterRequest) error {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return errors.New("กรุณากรอกให้ครบ")
	}
	RegisteruUser := registerDB{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := Register(RegisteruUser); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}