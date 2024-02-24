package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mail-client/internal/lib"
	"strconv"
	"strings"
)

func HandlerPopAuth(app *API) fiber.Handler {

	type request struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}

	return func(ctx *fiber.Ctx) error {

		var req request

		if err := ctx.BodyParser(&req); err != nil {
			return internal("Internal server error: cannot parse json")
		}

		if err := app.pop.Auth(req.User, req.Password); err != nil {
			//if 0 == strings.Compare(err.Error(), "failed at USER command: something went wrong: -ERR Unknown command: USER") {
			//	return bad("Already authenticated")
			//}
			//
			//if strings.Contains(err.Error(), "-ERR Invalid login or password") {
			//	return forbidden("Invalid credentials")
			//}

			return internal(fmt.Sprintf("Cannot auth: %s", err))
		}

		return respond(ctx, "successfully auth")
	}
}

func HandlerPopList(app *API) fiber.Handler {

	type response struct {
		Messages []*lib.Mail `json:"messages"`
	}

	return func(ctx *fiber.Ctx) error {

		var res response

		listedMessages, err := app.pop.ListAll()
		if err != nil {
			return internal(fmt.Sprintf("pop.ListAll: %s", err))
		}

		for _, listedMessage := range listedMessages {
			msg, err := app.pop.Retrieve(listedMessage.ID)
			if err != nil {
				return internal(fmt.Sprintf("pop.Retr: %s", err))
			}

			msg.Meta = listedMessage

			if len(msg.Body) > 128 {
				msg.Body = msg.Body[:128]
				msg.Body += "..."
			}

			res.Messages = append(res.Messages, msg)
		}

		return respond(ctx, res)
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
				return bad(fmt.Sprintf("mail id must be positive integer"))
			}
			return internal(fmt.Sprintf("%s [@HandlerPopRetrieve]", err.Error()))
		}

		msg, err := app.pop.Retrieve(id)
		if err != nil {
			if errors.Is(err, lib.ErrPop3Disconnected) {
				return forbidden(err.Error())
			}

			if strings.Contains(err.Error(), "There's no message") {
				return bad("unknown message")
			}

			return internal(fmt.Sprintf("%s [@HandlerPopRetrieve]", err.Error()))
		}

		res.Message = msg

		return respond(ctx, res)
	}
}

func HandlerPopStat(app *API) fiber.Handler {

	type response struct {
		Count int `json:"count"`
		Size  int `json:"size"`
	}

	return func(ctx *fiber.Ctx) error {

		var res response

		count, size, err := app.pop.Stat()
		if err != nil {
			if errors.Is(err, lib.ErrPop3Disconnected) {
				return forbidden(err.Error())
			}

			if strings.Contains(err.Error(), "There's no message") {
				return bad("unknown message")
			}

			return internal(fmt.Sprintf("%s [@HandlerPopRetrieve]", err.Error()))
		}

		res.Count = count
		res.Size = size

		return respond(ctx, res)
	}
}
