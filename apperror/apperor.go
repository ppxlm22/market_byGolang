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
//400 Bad Request
func NewBadRequest(message string) error {
	return &AppError{StatusCode: fiber.StatusBadRequest, Message: message}
}
//404 Not Found
func NewNotFound(message string) error {
	return &AppError{StatusCode: fiber.StatusNotFound, Message: message}
}
//500 Internal Server Error
func NewInternal(message string) error {
	return &AppError{StatusCode: fiber.StatusInternalServerError, Message: message}
}
//401 Unauthorized
func NewUnauthorized(message string) error {
	return &AppError{StatusCode: fiber.StatusUnauthorized, Message: message}
}
//409 Conflict
func NewConflict(message string) error {
	return &AppError{StatusCode: fiber.StatusConflict, Message: message}
}
//500 Internal Server Error
func NewInternalServerError(message string) error {
	return &AppError{StatusCode: fiber.StatusInternalServerError, Message: message}
}
//403 Forbidden
func NewForbidden(message string) error {
	return &AppError{StatusCode: fiber.StatusForbidden, Message: message}
}