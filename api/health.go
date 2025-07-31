package api

import "github.com/gofiber/fiber/v3"

func healthz(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
