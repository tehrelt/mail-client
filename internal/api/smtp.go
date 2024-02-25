package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mail-client/internal/dto"
	"mail-client/internal/lib"
	"os"
)

func HandlerSmtpSend(app *API) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var req dto.Message

		os.MkdirAll("./temp", 777)
		if len(req.Attachments) > 0 {
			for i, file := range req.Attachments {
				file, err := ctx.FormFile(file)
				if err != nil {
					return Internal(err.Error())
				}
				if err := ctx.SaveFile(file, fmt.Sprintf("./temp/%s", file.Filename)); err != nil {
					return Internal(err.Error())
				}
				req.Attachments[i] = fmt.Sprintf("./temp/%s", file.Filename)
			}
		}

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

		os.RemoveAll("./temp/")

		return Respond(ctx, "successfully sent")
	}
}
