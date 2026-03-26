package apperror

import (
    "errors"
    "log"
    "github.com/gofiber/fiber/v2"
)

func HandleError(c *fiber.Ctx, err error) error {
	var appErr *AppError

	if errors.As(err, &appErr) {
        return c.Status(appErr.StatusCode).JSON(appErr)
    }
	log.Printf("[SERVER ERROR]: %v", err)
	return c.Status(500).JSON(fiber.Map{"error": "เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์"})
}