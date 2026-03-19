package handler
import (
	"github.com/gofiber/fiber/v2"
	"go_shopmarket/login/service"
	"go_shopmarket/login/dto"

)

func NewHandler(s service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Login_Service(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "กรุณากรอกชื่อผู้ใช้และรหัสผ่าน",
		})
	}
	token, err := h.service.LoginUser(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	res := dto.LoginResponse{
		Token: token,
	}
	return c.Status(fiber.StatusOK).JSON(res)
}