package handler

import "github.com/gofiber/fiber/v2"

func HelloHandler(c *fiber.Ctx) error {
	return c.SendString("Hello, Handler!")
}
