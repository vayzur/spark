package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/vayzur/spark/xray"
	"github.com/xtls/xray-core/app/proxyman/command"
)

func addInbound(hsClient command.HandlerServiceClient) fiber.Handler {
	return func(c fiber.Ctx) error {
		b := c.Body()

		if err := xray.AddInbound(hsClient, b); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusCreated)
	}
}

func removeInbound(hsClient command.HandlerServiceClient) fiber.Handler {
	return func(c fiber.Ctx) error {
		tag := c.Params("tag")

		if tag != "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := xray.RemoveInbound(hsClient, tag); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusCreated)
	}
}
