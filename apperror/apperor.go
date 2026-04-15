package apperror
import (
	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	StatusCode int    `json:"code"`   
	Message    string `json:"message"`
}
func (e *AppError) Error() string {
	return e.Message
}
func NewBadRequest(message string) error {
	return &AppError{StatusCode: fiber.StatusBadRequest, Message: message}
}
func NewNotFound(message string) error {
	return &AppError{StatusCode: fiber.StatusNotFound, Message: message}
}
func NewInternal(message string) error {
	return &AppError{StatusCode: fiber.StatusInternalServerError, Message: message}
}
func NewUnauthorized(message string) error {
	return &AppError{StatusCode: fiber.StatusUnauthorized, Message: message}
}
func NewConflict(message string) error {
	return &AppError{StatusCode: fiber.StatusConflict, Message: message}
}
func NewInternalServerError(message string) error {
	return &AppError{StatusCode: fiber.StatusInternalServerError, Message: message}
}
