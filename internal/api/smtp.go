package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	libsmtp "mail-client/internal/lib"
	"net/smtp"
)

func HandlerSmtpAuth(app *API, client *smtp.Client) HandlerFunc {

	type request struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}

	return func(ctx *fiber.Ctx) error {

		var req request

		if err := ctx.BodyParser(&req); err != nil {
			return fiber.NewError(500, "Internal server error: cannot parse json")
		}

		if err := client.Auth(libsmtp.LoginAuth(req.User, req.Password)); err != nil {
			return fiber.NewError(500, fmt.Sprintf("Cannot auth: %s", err))
		}

		//if err := client.Verify(req.User); err != nil {
		//	return fiber.NewError(403, fmt.Sprintf("Forbidden: %s", err))
		//}

		return ctx.SendString("successfully auth")
	}
}
