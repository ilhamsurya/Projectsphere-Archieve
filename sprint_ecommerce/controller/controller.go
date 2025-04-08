package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ramadhan1445sprint/sprint_ecommerce/customErr"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
	"github.com/ramadhan1445sprint/sprint_ecommerce/svc"
)

type Controller struct {
	svc svc.SvcInterface
}

func NewController(svc svc.SvcInterface) *Controller {
	return &Controller{svc: svc}
}

func (c *Controller) CreateProduct(ctx *fiber.Ctx) error {
	product := &entity.Product{}

	// Check, if received JSON data is valid.
	if err := ctx.BodyParser(product); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": "internal server error"})
	}

	userId, _ := uuid.Parse(ctx.Locals("user_id").(string))
	product.UserID = userId

	err := c.svc.CreateProduct(*product)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "product added successfully",
	})
}

func (c *Controller) UpdateProduct(ctx *fiber.Ctx) error {
	product := &entity.Product{}

	id, err := uuid.Parse(ctx.Params("productId"))
	if err != nil {
		return err
	}

	userId, _ := uuid.Parse(ctx.Locals("user_id").(string))

	// Check, if received JSON data is valid.
	if err := ctx.BodyParser(product); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": "internal server error"})
	}

	product.ID = id
	product.UserID = userId

	err = c.svc.UpdateProduct(*product)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "product updated successfully",
	})
}

func (c *Controller) GetDetailProduct(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("productId"))
	if err != nil {
		return err
	}

	product, err := c.svc.GetDetailProduct(id)
	if err != nil {
		return customErr.NewNotFoundError("product not found")
	}

	productPayment, err := c.svc.GetProductSoldTotal(product.UserID)
	if err != nil {
		return customErr.NewInternalServerError("error query sold count")
	}

	bankAccount, _, err := c.svc.GetBankAccount(product.UserID.String())
	if err != nil {
		return customErr.NewInternalServerError(err.Error())
	}

	// Return status 200 OK.
	return ctx.JSON(fiber.Map{
		"message": "ok",
		"data": fiber.Map{
			"product": product,
			"seller": fiber.Map{
				"name":             productPayment.Name,
				"productSoldTotal": productPayment.TotalSold,
				"bankAccounts":     bankAccount,
			},
		},
	})
}

func (c *Controller) DeleteProduct(ctx *fiber.Ctx) error {
	userId, _ := uuid.Parse(ctx.Locals("user_id").(string))

	id, err := uuid.Parse(ctx.Params("productId"))
	if err != nil {
		return err
	}

	err = c.svc.DeleteProduct(id, userId)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "product deleted successfully",
	})
}

func (c *Controller) GetListProduct(ctx *fiber.Ctx) error {
	// Extract keys from the query parameters map
	keys := &entity.Key{}

	var products []entity.Product
	var limit, offset int = 0, 0
	var userId uuid.UUID

	ctx.QueryParser(keys)

	// check if userOnly == true
	if keys.UserOnly != nil && *keys.UserOnly {
		userId, _ = uuid.Parse(ctx.Locals("user_id").(string))
	}

	products, err := c.svc.GetListProduct(*keys, userId)
	if err != nil {
		return customErr.NewInternalServerError(err.Error())
	}

	if keys.Limit != nil && keys.Offset != nil {
		limit = *keys.Limit
		offset = *keys.Offset
	}

	// Return status 200 OK.
	return ctx.JSON(fiber.Map{
		"message": "ok",
		"data":    products,
		"meta": fiber.Map{
			"limit":  limit,
			"offset": offset,
			"total":  len(products),
		},
	})
}
