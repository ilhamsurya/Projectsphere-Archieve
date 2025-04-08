package utils

import (
	"errors"
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
)

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) map[string]string {
	// Define fields map.
	fields := map[string]string{}

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}

func ValidatePaymentRequest(paymentReq *entity.PaymentCreateRequest) (int, error) {
	if paymentReq.BankAccountID == nil {
		return 400, errors.New("bank account cant be null")
	}

	if paymentReq.PaymentProofImgUrl == nil {
		return 400, errors.New("payment proof image url cant be null")
	}

	if paymentReq.ProductID == nil {
		return 400, errors.New("product id cant be null")
	}

	if paymentReq.Quantity == nil {
		return 400, errors.New("quantity cant be null")
	}

	if *paymentReq.Quantity < 1 {
		return 400, errors.New("minimum quantity is 1")
	}

	if *paymentReq.BankAccountID == "" {
		return 400, errors.New("bank account cant be empty")
	}

	if reflect.TypeOf(*paymentReq.Quantity).Kind() != reflect.Int {
		return 400, errors.New("invalid quantity")
	}

	pattern := `^https://[^/]+\.s3\.[^/]+\.amazonaws\.com/.+\.(?:jpg|jpeg)$`

	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(*paymentReq.PaymentProofImgUrl) {
		return 400, errors.New("invalid image url")
	}

	return 200, nil
}

func ValidateCreateBankRequest(req *entity.BankAccountCreateRequest) (int, error) {
	if req.Name == nil {
		return 400, errors.New("bank name cant be null")
	}

	if req.AccountName == nil {
		return 400, errors.New("bank account name cant be null")
	}

	if req.AccountNumber == nil {
		return 400, errors.New("bank account number cant be null")
	}

	if len(*req.Name) < 5 || len(*req.Name) > 15 {
		return 400, errors.New("bank name should be minimum 5 and maximum 15 length")
	}

	if len(*req.AccountName) < 5 || len(*req.Name) > 15 {
		return 400, errors.New("bank account name should be minimum 5 and maximum 15 length")
	}

	if len(*req.AccountNumber) < 5 || len(*req.AccountNumber) > 15 {
		return 400, errors.New("bank account number should be minimum 5 and maximum 15 length")
	}

	return 200, nil
}

func ValidateUpdateBankRequest(req *entity.BankAccountUpdateRequest) (int, error) {
	if req.ID == nil {
		return 400, errors.New("bank id cant be null")
	}

	if *req.ID == "" {
		return 400, errors.New("invalid bank id")
	}

	if req.Name == nil {
		return 400, errors.New("bank name cant be null")
	}

	if req.AccountName == nil {
		return 400, errors.New("bank account name cant be null")
	}

	if req.AccountNumber == nil {
		return 400, errors.New("bank account number cant be null")
	}

	if len(*req.Name) < 5 || len(*req.Name) > 15 {
		return 400, errors.New("bank name should be minimum 5 and maximum 15 length")
	}

	if len(*req.AccountName) < 5 || len(*req.Name) > 15 {
		return 400, errors.New("bank account name should be minimum 5 and maximum 15 length")
	}

	if len(*req.AccountNumber) < 5 || len(*req.AccountNumber) > 15 {
		return 400, errors.New("bank account number should be minimum 5 and maximum 15 length")
	}

	return 200, nil
}

func ValidateStockUpdateRequest(req *entity.StockUpdateRequest) (int, error) {
	if req.Stock == nil {
		return 400, errors.New("stock cant be null")
	}

	if *req.Stock < 1 {
		return 400, errors.New("minimum stock is 0")
	}

	if req.ProductID == nil {
		return 400, errors.New("product id cant be null")
	}

	if *req.ProductID == "" {
		return 400, errors.New("product id cant be empty")
	}

	return 200, nil
}
