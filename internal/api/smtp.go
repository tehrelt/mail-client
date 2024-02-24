package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	libsmtp "mail-client/internal/lib"
	"net/smtp"
)

func HandlerSmtpAuth(app *API) HandlerFunc {

	type request struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}

	return func(ctx *fiber.Ctx) error {

		var req request

		if err := ctx.BodyParser(&req); err != nil {
			return fiber.NewError(500, "Internal server error: cannot parse json")
		}

		log.Debug(req)

		auth := libsmtp.LoginAuth(req.User, req.Password)
		if err := app.smtp.Auth(auth); err != nil {
			return fiber.NewError(500, fmt.Sprintf("Cannot auth: %s", err))
		}

		app.smtpAuth = auth

		return ctx.SendString("successfully auth")
	}
}

func HandlerSmtpSend(app *API) HandlerFunc {

	type request struct {
		From    string   `json:"from"`
		To      []string `json:"to"`
		Subject string   `json:"subject"`
		Body    string   `json:"body"`
	}

	return func(ctx *fiber.Ctx) error {

		var req request

		if err := ctx.BodyParser(&req); err != nil {
			return fiber.NewError(500, "Internal server error: cannot parse json")
		}

		log.Debug(req)

		if app.smtpAuth == nil {
			return fiber.NewError(403, "Authentication required")
		}

		//msg := []byte(
		//	"To: " + strings.Join(req.To, ",") + "\r\n" +
		//		"Subject: " + req.Subject + "\r\n" +
		//		"\r\n" +
		//		req.Body + "\r\n")

		msg := []byte(req.Body)

		log.Debug(string(msg))

		if err := smtp.SendMail(
			fmt.Sprintf("%s:%d", app.config.Smtp.Host, app.config.Smtp.Port),
			app.smtpAuth,
			req.From,
			req.To,
			msg,
		); err != nil {
			return fiber.NewError(500, fmt.Sprintf("Internal server error: %s", err))
		}

		return ctx.SendString("successfully sent")
	}
}
