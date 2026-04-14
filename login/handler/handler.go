package handler
import (
	"github.com/gofiber/fiber/v2"
	"go_shopmarket/login/service"
	"go_shopmarket/login/dto"
	"go_shopmarket/apperror"
	"github.com/go-playground/validator/v10"

)
var validate = validator.New()	

func NewHandler(s service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Login_Service(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequest("ข้อมูลที่ส่งมาไม่ถูกต้อง")
	}
	if err := validate.Struct(req); err != nil {
		return apperror.NewBadRequest(err.Error())
	}
	token, user, err := h.service.LoginUser(req)
	if err != nil {
		return err
	}
	res := dto.LoginResponse{
		Token: token,
		User:  user,
	}
	return c.Status(fiber.StatusOK).JSON(res)
}