package handlers

import "github.com/gofiber/fiber/v2"

func RespondWithError(code int, description string, c *fiber.Ctx) error {
	return c.Status(code).JSON(struct {
		Errors string `json:"errors"`
	}{Errors: description})
}
