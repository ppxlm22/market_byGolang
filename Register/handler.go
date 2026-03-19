package register
import (
	"github.com/gofiber/fiber/v2"
)
type Handler struct {
	service Service
}	

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Register_Service(c *fiber.Ctx) error {
	var req registerDB

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ข้อมูลไม่ถูกต้อง")
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "กรุณากรอกให้ครบ")
	}
	reqDB := registerDB{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	err := h.service.RegisterUser(reqDB)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "สมัครสมาชิกสำเร็จ"})

}	