package handlers

import (
	"github.com/gofiber/fiber/v3"
	"mail-client/internal/api"
	"net/smtp"
)

func HandlerSmtpAuth(app *api.API, client *smtp.Client) api.HandlerFunc {

	type request struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}

	return func(ctx fiber.Ctx) error {

		var req request

		if err := client.Auth(smtp.PlainAuth("", req.User, req.Password, app.Cfg.Smtp.Host)); err != nil {
			return fiber.NewError(500, "Cannot auth")
		}

		if err := client.Verify(req.User); err != nil {
			return fiber.NewError(403, "Forbidden: invalid credentials")
		}

		return ctx.SendString("successfully auth")
	}
}
