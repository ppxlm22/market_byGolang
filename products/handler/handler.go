package handler
import (
	"go_shopmarket/products/service"
	"go_shopmarket/products/dto"
	"github.com/gofiber/fiber/v2"
	"go_shopmarket/apperror"
	"github.com/go-playground/validator/v10"
)
var validate = validator.New()

func NewHandler(s service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	var req dto.Products
	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequest("ข้อมูลไม่ถูกต้อง")
	}
	if err := validate.Struct(req); err != nil {
		return apperror.NewBadRequest("ข้อมูลไม่ถูกต้อง")
	}
	if err := h.service.CreateProduct(req); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		
		"message": "สร้างสินค้าสำเร็จ",
	})			
	
}
func (h *Handler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.service.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(apperror.NewBadRequest(err.Error()))	
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": products,
	})
}
func (h *Handler) GetProductByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")	
	if err != nil {
		return apperror.NewBadRequest("ID ต้องเป็นตัวเลข")
	}
	product, err := h.service.GetProductByID(id)	
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"product": product,
	})
}	
func (h *Handler) UpdateProduct(c *fiber.Ctx) error {
	var req dto.Products
	id, err := c.ParamsInt("id")	
	if err != nil {
		return apperror.NewBadRequest("ID ต้องเป็นตัวเลข")
	}
	if err := c.BodyParser(&req); 
		err != nil {
			return apperror.NewBadRequest("ข้อมูลไม่ถูกต้อง")
	}
	if err := h.service.UpdateProduct(id, req);
		err != nil {
			return apperror.NewBadRequest("อัพเดตสินค้าไม่สำเร็จ")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "อัพเดตสินค้าสำเร็จ",
	})
}
func (h *Handler) DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.NewBadRequest("ID ต้องเป็นตัวเลข")
	}
	if err := h.service.DeleteProduct(id); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ลบสินค้าสำเร็จ",
	})
}
func (h *Handler) GetCategoryByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.NewBadRequest("ID ต้องเป็นตัวเลข")
	}
	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"category": category,
	})
}
func (h *Handler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		return apperror.NewInternalServerError("เกิดข้อผิดพลาดในการดึงข้อมูลหมวดหมู่")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"categories": categories,
	})
}
//
func (h *Handler) Checkout_service(c *fiber.Ctx) error {
	var req dto.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequest("ข้อมูลไม่ถูกต้อง")
	}
	if len(req.Items) == 0 {
		return apperror.NewBadRequest("ไม่มีสินค้าในคำสั่งซื้อ")
	}
	if err := h.service.Checkout(req); err != nil {
		return apperror.NewBadRequest("เกิดข้อผิดพลาดในการชำระเงิน")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ชำระเงินสำเร็จ",
	})
}