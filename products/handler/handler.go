package handler
import (
	"go_shopmarket/products/service"
	"go_shopmarket/products/dto"
	"github.com/gofiber/fiber/v2"
	"errors"
	"github.com/go-playground/validator/v10"
	"os"
	"log/slog"
)
var validate = validator.New()
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func NewHandler(s service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	var req dto.Products
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
		})
	}
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
		})
	}	
	if err := h.service.CreateProduct(req); err != nil {
		logger.Error("Failed to create product", 
			"error", err, 
			"product_name", req.Name,
			"product_price", req.Price,
			"product_category_id", req.CategoryID,
		)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "เกิดข้อผิดพลาดในการเพิ่มสินค้า",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "สร้างสินค้าสำเร็จ",
	})			
	
}
func (h *Handler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.service.GetAllProducts()
	if err != nil {
		logger.Error("Failed to get all products", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "เกิดข้อผิดพลาดในการดึงข้อมูลสินค้า",
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
		if errors.Is(err, service.ErrProductNotFound) {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "ไม่พบสินค้านี้ในระบบ",
        })
    }
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "เกิดข้อผิดพลาดในการดึงข้อมูลสินค้า",
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
				"error": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
		})
	}
	if err := h.service.UpdateProduct(id, req);
		err != nil {
			if errors.Is(err, service.ErrProductNotFound) {
        		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            		"error": "ไม่พบสินค้านี้ในระบบ",
        		})
    }
			logger.Error("Failed to update product")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "เกิดข้อผิดพลาดในการอัพเดตสินค้า",
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
		if errors.Is(err, service.ErrProductNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "ไม่พบสินค้านี้ในระบบ",
			})
    }
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "เกิดข้อผิดพลาดในการลบสินค้า",
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "เกิดข้อผิดพลาดในการดึงหมวดหมู่สินค้า",
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
			"error": "เกิดข้อผิดพลาดในการดึงหมวดหมู่สินค้า",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"categories": categories,
	})
}
//
func (h *Handler) Checkout_service(c *fiber.Ctx) error {
	var req dto.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
		})
	}
	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ต้องมีรายการสินค้าอย่างน้อยหนึ่งรายการ",
		})
	}
	if err := h.service.Checkout(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "เกิดข้อผิดพลาดในการชำระเงิน",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ชำระเงินสำเร็จ",
	})
}