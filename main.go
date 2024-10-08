package main

import (
	"time"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	app := fiber.New()

	app.Get("/", public)
	app.Post("/login", login)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("super_secret")},
	}))

	app.Get("/restricted", restricted)
	app.Get("/admin", admin)

	app.Listen(":3000")
}

func admin(c *fiber.Ctx) error {
	return c.SendString("Welcome Admin!")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims) //{"name": "jon", "admin": true}
	name := claims["name"].(string)
	return c.SendString("Welcome " + name + "\n")
}

func login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("super_secret"))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})

}

func public(c *fiber.Ctx) error {
	return c.SendString("THIS IS PUBLIC ROUTE")
}
