package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
	"github.com/ramadhan1445sprint/sprint_ecommerce/svc"
)

type stockController struct {
	svc svc.StockSvcInterface
}

func NewStockController(svc svc.StockSvcInterface) *stockController {
	return &stockController{svc: svc}
}

func (c *stockController) UpdateStock(ctx *fiber.Ctx) error {
	var message string
	
	productId := ctx.Params("productId")
	userId := ctx.Locals("user_id").(string)

	req := new(entity.StockUpdateRequest)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": "internal server error"})
	}

	req.ProductID = &productId

	status, err := c.svc.UpdateStock(req, productId, userId)

	if err != nil {
		message = err.Error()
	}else {
		message = "success"
	}

	resp := entity.Response{
		Message: message,
	}

	return ctx.Status(status).JSON(resp)
}
