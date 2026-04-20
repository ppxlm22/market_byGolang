package middleware

import (
	"log/slog"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")	
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "ไม่อนุญาตให้เข้าถึง",
			})
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _,ok := token.Method.(*jwt.SigningMethodHMAC); ok{
				return nil, fiber.ErrUnauthorized
			}
			secret := os.Getenv("JWT_SECRET")
			return []byte(secret), nil
		})
		if 	err != nil || !token.Valid {
			slog.Warn("Unanthorized access attempt","error", err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token หมดอายุ",
			})
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("role", claims["role"],)
			c.Locals("user_id",claims["user_id"])
		}
		return c.Next()
	}
}
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		user_id := c.Locals("user_id")
		roleStr, ok := role.(string)
		//string and role Admin
		if !ok || roleStr != "admin" {
			slog.Warn("Forbidden access attempt: Not an admin",
				"role_tried", roleStr,
				"userID", user_id,
			)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "สิทธิ์การเข้าถึงถูกปฏิเสธ",
			})
		}
		return c.Next()
	}
}