package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_ecommerce/customErr"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if customErr, ok := err.(customErr.CustomError); ok {
		code = customErr.Status()
		return ctx.Status(code).JSON(fiber.Map{
			"message": customErr.Error(),
		})
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": err.Error(),
	})
}
