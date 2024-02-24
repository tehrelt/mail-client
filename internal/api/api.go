package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"mail-client/internal/config"
	"mail-client/internal/lib"
	"net/smtp"
)

type HandlerFunc func(ctx *fiber.Ctx) error

type API struct {
	app      *fiber.App
	smtp     *smtp.Client
	smtpAuth smtp.Auth
	pop      *lib.Pop3

	config *config.AppConfig
}

func Start(cfg *config.AppConfig) error {
	app := fiber.New()

	smtp, err := smtp.Dial(fmt.Sprintf("%s:%d", cfg.Smtp.Host, cfg.Smtp.Port))
	if err != nil {
		return err
	}

	pop3 := lib.NewPop(cfg.Pop3)

	api := &API{
		app: app,

		smtp:     smtp,
		smtpAuth: nil,

		pop: pop3,

		config: cfg,
	}

	api.app.Use(logger.New())
	api.configure()

	return api.app.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func (api *API) configure() {
	api.app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.JSON(&fiber.Map{
			"message": "pong",
		})
	})

	smtp := api.app.Group("/smtp")
	smtp.Post("/auth", HandlerSmtpAuth(api))
	smtp.Post("/send", HandlerSmtpSend(api))

	pop3 := api.app.Group("/pop3")
	pop3.Post("/auth", HandlerPopAuth(api))
	pop3.Get("/list", HandlerPopList(api))
}
