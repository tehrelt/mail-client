package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mail-client/internal/dto"
)

func HandlerSmtpAuth(app *API) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var req dto.User

		if err := ctx.BodyParser(&req); err != nil {
			return fiber.NewError(500, "Internal server error: cannot parse json")
		}

		if err := app.smtp.Auth(&req); err != nil {
			return internal(fmt.Sprintf("smtp.Auth: %s", err.Error()))
		}

		return respond(ctx, "successfully authenticated")
	}
}

func HandlerSmtpSend(app *API) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var req dto.Message

		if err := ctx.BodyParser(&req); err != nil {
			return fiber.NewError(500, "Internal server error: cannot parse json")
		}

		if err := app.smtp.Send(&req); err != nil {
			return internal(fmt.Sprintf("smtp.Send: %s", err.Error()))
		}

		return respond(ctx, "successfully sent")
	}
}
