package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/vayzur/spark/config"
	"github.com/vayzur/spark/internal/auth"
	"github.com/xtls/xray-core/app/proxyman/command"
)

func authMiddleware(c fiber.Ctx) error {
	h := c.Get("Authorization")
	if h == "" {
		return fiber.ErrUnauthorized
	}

	if err := auth.VerifyRollingHash(h); err != nil {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}

func requireJSON(c fiber.Ctx) error {
	ct := c.Get(fiber.HeaderContentType)
	if ct != fiber.MIMEApplicationJSON {
		return c.Status(fiber.StatusUnsupportedMediaType).
			JSON(fiber.Map{"error": "Content-Type must be application/json"})
	}
	return c.Next()
}

func NewAPIServer(hsClient command.HandlerServiceClient) *fiber.App {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
	})

	app.Use(authMiddleware)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/healthz", healthz)

	inbounds := app.Group("/inbounds", requireJSON)
	inbounds.Post("/", addInbound(hsClient))
	inbounds.Delete("/:tag", removeInbound(hsClient))

	return app
}

func StartAPIServerTLS(addr string, app *fiber.App) error {
	return app.Listen(addr, fiber.ListenConfig{
		CertFile:      config.AppConfig.TLS.CertFile,
		CertKeyFile:   config.AppConfig.TLS.KeyFile,
		EnablePrefork: config.AppConfig.Server.Prefork,
	})
}

func StartAPIServer(addr string, app *fiber.App) error {
	return app.Listen(addr, fiber.ListenConfig{
		EnablePrefork: config.AppConfig.Server.Prefork,
	})
}
