package middleware

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type CORS struct {
	MW fiber.Handler
}

func NewCORSMiddleware() *CORS {
	return &CORS{MW: func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Methods", "POST,PUT,DELETE,GET")
		c.Set("Access-Control-Allow-Headers", "Content-Type")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Origin", c.Get("Origin"))
		if c.Method() == http.MethodOptions {
			return nil
		}

		return c.Next()
	},
	}
}
