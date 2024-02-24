package api

import (
	"fmt"
	"github.com/bytbox/go-pop3"
	"github.com/gofiber/fiber/v3"
	"mail-client/internal/config"
	"net/smtp"
)

type API struct {
	*fiber.App
	smtp *smtp.Client
	pop  *pop3.Client
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
		App:  app,
		smtp: smtp,
		pop:  pop3,
	}

	api.configure()

	return api.Listen(":7000")
}

func (api *API) configure() {
	api.Get("/ping", func(ctx fiber.Ctx) error {
		return ctx.JSON(&fiber.Map{
			"message": "pong",
		})
	})
}
