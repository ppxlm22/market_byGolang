package handler

import (
	"errors"
	"go_shopmarket/register/dto"
	"go_shopmarket/register/service"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func NewHandler(s service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Register_Service(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
		})
	}

	// check Validation
	if err := validate.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		
		if errors.As(err, &validationErrors) {
			firstErr := validationErrors[0] 
			var customMessage string

			switch firstErr.Tag() {
			case "required":
				customMessage = "กรุณากรอก " + firstErr.Field() + " ให้ครบถ้วน"
			case "email":
				customMessage = "รูปแบบอีเมลไม่ถูกต้อง"
			case "min":
				customMessage = firstErr.Field() + " สั้นเกินไป ต้องมีอย่างน้อย " + firstErr.Param() + " ตัวอักษร"
			default:
				customMessage = "ข้อมูลไม่ถูกต้องตามเงื่อนไข"
			}

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": customMessage,
			})
		}
		
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลไม่ถูกต้อง",
		})
	}

	Userrespone, err := h.service.RegisterUser(req)
	if err != nil {
		if errors.Is(err, service.ErrUserDuplicate) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Username หรือ Email นี้มีผู้ใช้งานแล้ว",
			})
		}
		
		logger.Error("Failed to register user",
			"error", err,
			"username", req.Username,
			"email", req.Email,
		)
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "เกิดข้อผิดพลาดในการลงทะเบียน กรุณาลองใหม่อีกครั้ง",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(Userrespone)
}