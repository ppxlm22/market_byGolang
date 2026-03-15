package shopHandler

import (
	"database/sql"
	"go_shopmarket/database"
	_ "go_shopmarket/database"
	"go_shopmarket/models"
	"strconv"
	"strings"

	"fmt"
	"time"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)
func AuthRequired(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(401).JSON(fiber.Map{"error": "ไม่อนุญาตให้เข้าถึง: กรุณา Login ก่อน"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil 
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "ไม่อนุญาตให้เข้าถึง: Token ไม่ถูกต้องหรือหมดอายุ"})
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Locals("username", claims["username"]) 
		c.Locals("role", claims["role"])         
	}

	return c.Next()
}
func AdminRequired(c *fiber.Ctx) error {
	role := c.Locals("role")

	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden: คุณไม่มีสิทธิ์ใช้งานส่วนนี้ (ต้องเป็น Admin เท่านั้น)",
		})
	}
	return c.Next()
}

func Register(c *fiber.Ctx) error {
	type RegisterRequest struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
	var req RegisterRequest
	if err := c.BodyParser(&req); 
	err != nil {
		return c.Status(400).JSON(fiber.Map{"ERROR":"Data fail"})
	} 
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"ERROR": "Hash fail"})
	}
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`
	_, err = database.DB.Exec(query, req.Username, req.Email, string(hashedPassword))

	if err != nil {
        return c.Status(500).JSON(fiber.Map{"ERROR_LOG": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Register Success"})
}	
func Login(c *fiber.Ctx) error {
    type Loginrequest struct{
        Username string `json:"username"`
        Password string `json:"password"`
    }

    var req Loginrequest
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"ERROR": "Data fail"})
    }
    fmt.Println("Postman: [", req.Username, "]")

    var hashedPassword string
    query := `SELECT password_hash FROM users WHERE username = $1`
    err := database.DB.QueryRow(query, req.Username).Scan(&hashedPassword)
    
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "NOT Found User"})
    }
    
    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "รหัสผ่านไม่ถูกต้อง"})
    }
    
    claims := jwt.MapClaims{
        "username": req.Username,
        "role":     "customer",
        "exp":      time.Now().Add(time.Hour * 72).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, err := token.SignedString([]byte("secret")) 
    
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "ไม่สามารถสร้าง Token ได้"})
    }

    return c.JSON(fiber.Map{
        "message": "Login Success!",
        "user":    req.Username,
        "token":   t, 
    })
}


func Getproduct(c *fiber.Ctx) error {
	id := c.Params("id")
	fmt.Println(id)

	var p models.Products

	query := `SELECT id, name, price, stock, category_id, create_at FROM public.products WHERE id = $1;`

	err := database.DB.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.Create_at)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString("Product not found")
		}
		return err
	}

	return c.JSON(&p)
}
func GetAllproducts(c *fiber.Ctx) error {
	var p models.Products
	var proDucts []models.Products

	query := `SELECT * FROM public.products;`
	rows, _ := database.DB.Query(query)

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.Create_at); err != nil {
			return err
		}
		fmt.Println(p)
		proDucts = append(proDucts, p)
	}

	return c.JSON(proDucts)

}
func CreateProduct(c *fiber.Ctx) error {
	fmt.Println("CreateProduct")
	p := new(models.Products)
	if err := c.BodyParser(p); err != nil {
		return err
	}
	_, err := database.DB.Exec("INSERT INTO public.products (name, price, stock, category_id ) VALUES ($1, $2, $3, $4)", p.Name, p.Price, p.Stock, p.CategoryID)
	if err != nil {
		return err
	}
	return c.JSON(p)

}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	p := new(models.Products)
	if err := c.BodyParser(p); err != nil {
		return err
	}
	_, err := database.DB.Exec("UPDATE public.products SET name = $1, price = $2,stock = $3 ,category_id = $4 WHERE id = $5", p.Name, p.Price, p.Stock, p.CategoryID, id)
	if err != nil {
		return err
	}
	p.ID, _ = strconv.Atoi(id)
	return c.JSON(p)
}
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := database.DB.Exec("DELETE FROM public.products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
