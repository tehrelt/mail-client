package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mail-client/internal/dto"
	"mail-client/internal/lib"
	"os"
	"strings"
)

func HandlerSmtpSend(app *API) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		form, err := ctx.MultipartForm()
		if err != nil {
			return Internal(fmt.Sprintf("ctx.MultipartForm: %s", err.Error()))
		}

		from := form.Value["from"][0]
		to := form.Value["to"][0]
		subject := form.Value["subject"][0]
		body := form.Value["body"][0]

		os.MkdirAll("./temp", 0777)
		var attachments []string
		files := form.File
		for _, file := range files {
			dest := fmt.Sprintf("./temp/%s", file[0].Filename)
			err := ctx.SaveFile(file[0], dest)
			if err != nil {
				return err
			}
			attachments = append(attachments, dest)
		}

		if ctx.Locals("connection") == nil {
			return Internal(fmt.Sprintf("smtp.Send: where connection"))
		}

		connection := ctx.Locals("connection").(*lib.Smtp)
		if err := connection.SendMessage(&dto.Message{
			From:        from,
			To:          strings.Split(to, ","),
			Subject:     subject,
			Body:        body,
			Attachments: attachments,
		}); err != nil {
			return Internal(fmt.Sprintf("smtp.Send: %s", err.Error()))
		}

		os.RemoveAll("./temp/")

		return Respond(ctx, "successfully sent")
	}
}
