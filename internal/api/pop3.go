package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mail-client/internal/lib"
	"strconv"
	"strings"
)

func HandlerPopList(app *API) fiber.Handler {

	type response struct {
		Messages []*lib.Mail `json:"messages"`
	}

	return func(ctx *fiber.Ctx) error {

		var res response

		if ctx.Locals("connection") == nil {
			return Internal(fmt.Sprintf("pop.ListAll: where connection"))
		}

		connection := ctx.Locals("connection").(*lib.Pop3)

		listedMessages, err := connection.ListAll()
		if err != nil {
			return Internal(fmt.Sprintf("pop.ListAll: %s", err))
		}

		for _, listedMessage := range listedMessages {
			msg, err := connection.Retrieve(listedMessage.ID)
			if err != nil {
				return Internal(fmt.Sprintf("pop.Retr: %s", err))
			}

			msg.Meta = listedMessage

			if len(msg.Body) > 128 {
				msg.Body = msg.Body[:128]
				msg.Body += "..."
			}

			res.Messages = append(res.Messages, msg)
		}

		return Respond(ctx, res)
	}
}

func HandlerPopRetrieve(app *API) fiber.Handler {

	type response struct {
		Message *lib.Mail `json:"message"`
	}

	return func(ctx *fiber.Ctx) error {

		var res response

		idParam := ctx.Params("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			if errors.Is(err, strconv.ErrSyntax) {
				return Bad(fmt.Sprintf("mail id must be positive integer"))
			}
			return Internal(fmt.Sprintf("%s [@HandlerPopRetrieve]", err.Error()))
		}

		if ctx.Locals("connection") == nil {
			return Internal(fmt.Sprintf("pop.ListAll: where connection"))
		}

		connection := ctx.Locals("connection").(*lib.Pop3)

		msg, err := connection.Retrieve(id)
		if err != nil {
			if errors.Is(err, lib.ErrPop3Disconnected) {
				return Forbidden(err.Error())
			}

			if strings.Contains(err.Error(), "There's no message") {
				return Bad("unknown message")
			}

			return Internal(fmt.Sprintf("%s [@HandlerPopRetrieve]", err.Error()))
		}

		res.Message = msg

		return Respond(ctx, res)
	}
}

func HandlerPopStat(app *API) fiber.Handler {

	type response struct {
		Count int `json:"count"`
		Size  int `json:"size"`
	}

	return func(ctx *fiber.Ctx) error {

		var res response
		if ctx.Locals("connection") == nil {
			return Internal(fmt.Sprintf("pop.ListAll: where connection"))
		}

		connection := ctx.Locals("connection").(*lib.Pop3)
		count, size, err := connection.Stat()
		if err != nil {
			if errors.Is(err, lib.ErrPop3Disconnected) {
				return Forbidden(err.Error())
			}

			if strings.Contains(err.Error(), "There's no message") {
				return Bad("unknown message")
			}

			return Internal(fmt.Sprintf("%s [@HandlerPopRetrieve]", err.Error()))
		}

		res.Count = count
		res.Size = size

		return Respond(ctx, res)
	}
}
