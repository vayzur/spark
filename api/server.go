package api

import (
	"crypto/tls"

	"github.com/gofiber/fiber/v3"
	"github.com/vayzur/spark/config"
	"github.com/vayzur/spark/internal/auth"
	"github.com/xtls/xray-core/app/proxyman/command"
	"golang.org/x/crypto/acme/autocert"
)

func authMiddleware(c fiber.Ctx) error {
	h := c.Get("Authorization") // "<ts>:<hex-sig>"
	if h == "" {
		return fiber.ErrUnauthorized
	}

	if err := auth.Verify(h); err != nil {
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

	// app.Use(authMiddleware)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/healthz", healthz)

	inbounds := app.Group("/inbounds", requireJSON)
	inbounds.Post("/", addInbound(hsClient))
	inbounds.Delete("/:tag", removeInbound(hsClient))

	return app
}

func StartAPIServerTLS(addr string, app *fiber.App) {
	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.AppConfig.Domain),
		Cache:      autocert.DirCache("./certs"),
	}

	app.Listen(config.AppConfig.APIServerAddr, fiber.ListenConfig{
		AutoCertManager: certManager,
		TLSMinVersion:   tls.VersionTLS12,
		EnablePrefork:   true,
	})
}

func StartAPIServer(addr string, app *fiber.App) {
	app.Listen(config.AppConfig.APIServerAddr, fiber.ListenConfig{
		EnablePrefork: true,
	})
}
