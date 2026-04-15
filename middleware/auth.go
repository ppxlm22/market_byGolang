package middleware

import (
	"strings"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go_shopmarket/apperror"
)
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")	
		if authHeader == "" {
			return apperror.NewUnauthorized("ไม่อนุญาติให้เข้าถึง")
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			secret := os.Getenv("JWT_SECRET")
			return []byte(secret), nil
		})
		if 	err != nil || !token.Valid {
			return apperror.NewUnauthorized("Token ไม่ถูกต้องหรือหมดอายุ")
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("role", claims["role"])
		}
		return c.Next()
	}
}
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")

		if role != "admin" {
			return apperror.NewForbidden("สิทธิ์การเข้าถึงถูกปฏิเสธ")
		}
		return c.Next()
	}
}