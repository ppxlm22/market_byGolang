package handler
import (
	"github.com/gofiber/fiber/v2"
	"go_shopmarket/login/service"
	"go_shopmarket/login/dto"
	"github.com/go-playground/validator/v10"
	"errors"

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
		})
	}
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลไม่ถูกต้อง",
		})
	}
	token, user, err := h.service.LoginUser(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "เกิดข้อผิดพลาดภายในระบบ",
		})
	}
	res := dto.LoginResponse{
		Token: token,
		User:  user,
	}
	return c.Status(fiber.StatusOK).JSON(res)
}