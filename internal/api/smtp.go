package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mail-client/internal/dto"
	"mail-client/internal/lib"
)

func HandlerSmtpSend(app *API) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var req dto.Message

		if err := ctx.BodyParser(&req); err != nil {
			return fiber.NewError(500, "Internal server error: cannot parse json")
		}

		if ctx.Locals("connection") == nil {
			return Internal(fmt.Sprintf("pop.ListAll: where connection"))
		}

		connection := ctx.Locals("connection").(*lib.Smtp)

		if err := connection.SendMessage(&req); err != nil {
			return Internal(fmt.Sprintf("smtp.Send: %s", err.Error()))
		}

		return Respond(ctx, "successfully sent")
	}
}
