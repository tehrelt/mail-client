package api

import (
	"fmt"
	"github.com/bytbox/go-pop3"
	"github.com/gofiber/fiber/v2"
	"mail-client/internal/config"
	"net/smtp"
)

type HandlerFunc func(ctx *fiber.Ctx) error

type API struct {
	app  *fiber.App
	smtp *smtp.Client
	pop  *pop3.Client

	Cfg *config.AppConfig
}

func Start(cfg *config.AppConfig) error {
	app := fiber.New()

	smtp, err := smtp.Dial(fmt.Sprintf("%s:%d", cfg.Smtp.Host, cfg.Smtp.Port))
	if err != nil {
		return err
	}

	pop3, err := pop3.Dial(fmt.Sprintf("%s:%d", cfg.Pop3.Host, cfg.Pop3.Port))
	if err != nil {
		return err
	}

	api := &API{
		app:  app,
		smtp: smtp,
		pop:  pop3,
		Cfg:  cfg,
	}

	api.configure()

	return api.app.Listen(":7000")
}

func (api *API) configure() {
	api.app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.JSON(&fiber.Map{
			"message": "pong",
		})
	})

	smtp := api.app.Group("/smtp")
	smtp.Post("/auth", HandlerSmtpAuth(api, api.smtp))

	//pop3 := api.app.Group("/pop3")
}
