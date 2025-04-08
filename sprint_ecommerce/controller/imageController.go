package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_ecommerce/customErr"
	"github.com/ramadhan1445sprint/sprint_ecommerce/svc"
)

type ImageController struct {
	svc svc.ImageSvc
}

func NewImageController(svc svc.ImageSvc) *ImageController {
	return &ImageController{
		svc: svc,
	}
}

func (c *ImageController) UploadImage(ctx *fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return customErr.NewInternalServerError("failed to retrieve file")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return customErr.NewInternalServerError("failed to open file")
	}

	url, err := c.svc.UploadImage(file, fileHeader)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"imageUrl": url,
	})

}
