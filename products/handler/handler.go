package handler
import (
	"go_shopmarket/products/service"
	"go_shopmarket/products/dto"
	"github.com/gofiber/fiber/v2"
)
func NewHandler(s service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	var req dto.Products
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลไม่ถูกต้อง",
		})
	}
	if err := h.service.CreateProduct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "สร้างสินค้าสำเร็จ",
	})
}
