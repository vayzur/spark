package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/vayzur/spark/internal/auth"
	"github.com/vayzur/spark/internal/config"
	"github.com/vayzur/spark/pkg/client/xray"
)

type Server struct {
	addr       string
	app        *fiber.App
	XrayClient *xray.XrayClient
}

func NewServer(addr string, xrayClient *xray.XrayClient) *Server {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
	})
	s := &Server{
		addr:       addr,
		app:        app,
		XrayClient: xrayClient,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// s.app.Use(authMiddleware)

	s.app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	s.app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())

	api := s.app.Group("/api")
	v1 := api.Group("/v1")

	inbounds := v1.Group("/inbounds")
	inbounds.Post("", requireJSON, s.AddInbound)
	inbounds.Delete("/:tag", s.RemoveInbound)
}

func (s *Server) StartTLS() error {
	return s.app.Listen(s.addr, fiber.ListenConfig{
		CertFile:      config.AppConfig.TLS.CertFile,
		CertKeyFile:   config.AppConfig.TLS.KeyFile,
		EnablePrefork: config.AppConfig.Prefork,
	})
}

func (s *Server) Start() error {
	return s.app.Listen(s.addr, fiber.ListenConfig{
		EnablePrefork: config.AppConfig.Prefork,
	})
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}

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
