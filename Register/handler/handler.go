package handler
import (
	"github.com/gofiber/fiber/v2"
	"go_shopmarket/register/service"
	"go_shopmarket/register/dto"
	"go_shopmarket/apperror"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func NewHandler(s service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Register_Service(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequest("ข้อมูลที่ส่งมาไม่ถูกต้อง")
	}
	if err := validate.Struct(req); err != nil {
		return apperror.NewBadRequest(err.Error())
	}
	reqDB := dto.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	Userrespone, err := h.service.RegisterUser(reqDB)
	if err != nil {
		return apperror.NewInternalServerError("เกิดข้อผิดพลาดในการลงทะเบียน")
	}

	return c.Status(fiber.StatusCreated).JSON(Userrespone)

}	