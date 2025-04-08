package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ramadhan1445sprint/sprint_ecommerce/crypto"
	"github.com/ramadhan1445sprint/sprint_ecommerce/customErr"
)

func Authorization(ctx *fiber.Ctx) error {
	auth := ctx.Get("Authorization")
	if auth == "" {
		return customErr.NewUnauthorizedError("token not found")
	}

	splitted := strings.Split(auth, " ")

	if splitted[0] != "Bearer" {
		return customErr.NewUnauthorizedError("invalid token")
	}

	payload, err := crypto.VerifyToken(splitted[1])
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return customErr.NewUnauthorizedError("token expired")
		}
		return customErr.NewUnauthorizedError(err.Error())
	}

	ctx.Locals("user_id", payload.Id)
	ctx.Locals("username", payload.Username)
	ctx.Locals("name", payload.Name)
	return ctx.Next()
}

func ProductPageAuth(ctx *fiber.Ctx) error {
	userOnly := ctx.Query("userOnly")

	if userOnly == "true" {
		auth := ctx.Get("Authorization")
		if auth == "" {
			return customErr.NewUnauthorizedError("token not found")
		}

		splitted := strings.Split(auth, " ")

		if splitted[0] != "Bearer" {
			return customErr.NewUnauthorizedError("invalid token")
		}

		payload, err := crypto.VerifyToken(splitted[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return customErr.NewUnauthorizedError("token expired")
			}
			return customErr.NewUnauthorizedError(err.Error())
		}

		ctx.Locals("user_id", payload.Id)
	}

	return ctx.Next()
}
