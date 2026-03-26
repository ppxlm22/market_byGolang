package apperror

import (
	"github.com/gofiber/fiber/v2"
)
type AppError struct {
	StatusCode int    `json:"-"`   
    Message    string `json:"error"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrNotFound = &AppError{
		StatusCode: fiber.StatusNotFound, 
		Message:    "ไม่พบข้อมูลที่ต้องการ",
	}

	ErrUnauthorized = &AppError{
		StatusCode: fiber.StatusUnauthorized, 
		Message:    "กรุณาเข้าสู่ระบบก่อนทำรายการ",
	}

	ErrForbidden = &AppError{
		StatusCode: fiber.StatusForbidden, 
		Message:    "คุณไม่มีสิทธิ์เข้าถึงส่วนนี้",
	}

	ErrBadRequest = &AppError{
		StatusCode: fiber.StatusBadRequest, 
		Message:    "ข้อมูลที่ส่งมาไม่ถูกต้อง",
	}
    
    ErrOutOfStock = &AppError{
		StatusCode: fiber.StatusBadRequest, 
		Message:    "สินค้าในสต็อกไม่เพียงพอ",
	}
)


