package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_ecommerce/svc"
)

type healthController struct {
	svc svc.HealthSvcInterface
}

func NewHealthController(svc svc.HealthSvcInterface) *healthController {
	return &healthController{svc: svc}
}

func (c *healthController) Get(ctx *fiber.Ctx) error {
	status, err := c.svc.GetStatus()
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"status":  status,
	})
}
