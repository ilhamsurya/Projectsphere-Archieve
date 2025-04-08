package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
	"github.com/ramadhan1445sprint/sprint_ecommerce/svc"
)

type UserController struct {
	svc svc.UserSvc
}

func NewUserController(svc svc.UserSvc) *UserController {
	return &UserController{svc: svc}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	var user entity.User
	err := ctx.BodyParser(&user)
	if err != nil {
		return err
	}

	token, err := c.svc.RegisterUser(&user)
	if err != nil {
		return err
	}

	return ctx.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": fiber.Map{
			"username":    user.Username,
			"name":        user.Name,
			"accessToken": token,
		},
	})

}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var creds entity.Credential
	err := ctx.BodyParser(&creds)
	if err != nil {
		return err
	}

	user, token, err := c.svc.Login(creds)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "User logged in successfully",
		"data": fiber.Map{
			"username":    user.Username,
			"name":        user.Name,
			"accessToken": token,
		},
	})
}
