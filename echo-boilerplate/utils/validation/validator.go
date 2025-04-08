package validation

import (
	"reflect"
	"strings"

	"github.com/JesseNicholas00/HaloSuster/utils/validation/image"
	"github.com/JesseNicholas00/HaloSuster/utils/validation/intlen"
	"github.com/JesseNicholas00/HaloSuster/utils/validation/iso8601"
	"github.com/JesseNicholas00/HaloSuster/utils/validation/nip"
	"github.com/JesseNicholas00/HaloSuster/utils/validation/phone"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type EchoValidator struct {
	validator *validator.Validate
}

func (e *EchoValidator) Validate(i interface{}) error {
	return e.validator.Struct(i)
}

var customFields = []customField{
	{
		Tag:       "phoneNum",
		Validator: phone.ValidatePhoneNumber,
	},
	{
		Tag:       "nip",
		Validator: nip.ValidateNip,
	},
	{
		Tag:       "nipNurse",
		Validator: nip.ValidateNipNurse,
	},
	{
		Tag:       "nipIt",
		Validator: nip.ValidateNipIt,
	},
	{
		Tag:       "imageExt",
		Validator: image.ValidateImageExtension,
	},
	{
		Tag:       "iso8601",
		Validator: iso8601.ValidateIso8601,
	},
	{
		Tag:       "intlen",
		Validator: intlen.ValidateIntLen,
	},
}

type customField struct {
	Tag       string
	Validator validator.Func
}

func NewEchoValidator() echo.Validator {
	validator := validator.New(validator.WithRequiredStructEnabled())

	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	for _, customField := range customFields {
		validator.RegisterValidation(customField.Tag, customField.Validator)
	}

	return &EchoValidator{
		validator: validator,
	}
}
