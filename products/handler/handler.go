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
func (h *Handler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.service.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": products,
	})
}
func (h *Handler) GetProductByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")	
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID ต้องเป็นตัวเลข",
		})
	}
	product, err := h.service.GetProductByID(id)	
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"product": product,
	})
}	
func (h *Handler) UpdateProduct(c *fiber.Ctx) error {
	var req dto.Products
	id, err := c.ParamsInt("id")	
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID ต้องเป็นตัวเลข",
		})
	}
	if err := c.BodyParser(&req); 
		err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ข้อมูลไม่ถูกต้อง",
			})
	}
	if err := h.service.UpdateProduct(id, req);
		err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "อัพเดตสินค้าสำเร็จ",
	})
}
func (h *Handler) DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID ต้องเป็นตัวเลข",
		})
	}
	if err := h.service.DeleteProduct(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ลบสินค้าสำเร็จ",
	})
}
func (h *Handler) GetCategoryByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID ต้องเป็นตัวเลข",
		})
	}
	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"category": category,
	})
}
func (h *Handler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"categories": categories,
	})
}
func (h *Handler) Checkout_service(c *fiber.Ctx) error {
	var req dto.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลไม่ถูกต้อง",
		})
	}
	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ไม่มีสินค้าในตะกร้า",
		})
	}
	if err := h.service.Checkout(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ชำระเงินสำเร็จ",
	})
}