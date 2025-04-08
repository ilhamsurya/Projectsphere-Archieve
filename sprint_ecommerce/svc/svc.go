package svc

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ramadhan1445sprint/sprint_ecommerce/customErr"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
	"github.com/ramadhan1445sprint/sprint_ecommerce/repo"
)

type SvcInterface interface {
	CreateProduct(product entity.Product) error
	GetDetailProduct(id uuid.UUID) (entity.Product, error)
	UpdateProduct(product entity.Product) error
	DeleteProduct(id uuid.UUID, userId uuid.UUID) error
	GetListProduct(keys entity.Key, userId uuid.UUID) ([]entity.Product, error)
	GetProductSoldTotal(userId uuid.UUID) (entity.ProductPayment, error)
	GetBankAccount(userId string) ([]entity.BankAccountGetResponse, int, error)
}

func NewSvc(repo repo.RepoInterface) SvcInterface {
	return &svc{repo: repo}
}

type svc struct {
	repo repo.RepoInterface
}

func ValidateCondition(fl validator.FieldLevel) bool {
	condition := fl.Field().String()
	return condition == "new" || condition == "second"
}

func (s *svc) CreateProduct(product entity.Product) error {
	validate := validator.New()
	validate.RegisterValidation("validCondition", ValidateCondition)

	if err := validate.Struct(product); err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	err := s.repo.CreateProduct(product)
	if err != nil {
		return customErr.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *svc) GetDetailProduct(id uuid.UUID) (entity.Product, error) {
	product, err := s.repo.GetDetailProduct(id)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *svc) UpdateProduct(product entity.Product) error {
	validate := validator.New()
	validate.RegisterValidation("validCondition", ValidateCondition)

	if err := validate.Struct(product); err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	if err := s.repo.UpdateProduct(product); err != nil {
		if err.Error() == "product not found" {
			err = customErr.NewNotFoundError(err.Error())
		} else if err.Error() == "product is not owned by this user" {
			err = customErr.NewForbiddenError(err.Error())
		} else {
			err = customErr.NewInternalServerError(err.Error())
		}
		return err
	}

	return nil
}

func (s *svc) DeleteProduct(id uuid.UUID, userId uuid.UUID) error {
	if err := s.repo.DeleteProduct(id, userId); err != nil {
		if err.Error() == "product not found" {
			err = customErr.NewNotFoundError(err.Error())
		} else if err.Error() == "product is not owned by this user" {
			err = customErr.NewForbiddenError(err.Error())
		} else {
			err = customErr.NewInternalServerError(err.Error())
		}
		return err
	}

	return nil
}

func (s *svc) GetListProduct(keys entity.Key, userId uuid.UUID) ([]entity.Product, error) {
	product, err := s.repo.GetListProduct(keys, userId)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *svc) GetProductSoldTotal(userId uuid.UUID) (entity.ProductPayment, error) {
	productPayment, err := s.repo.GetProductSoldTotal(userId)
	if err != nil {
		return productPayment, err
	}

	return productPayment, nil
}

func (s *svc) GetBankAccount(userId string) ([]entity.BankAccountGetResponse, int, error) {
	res, err := s.repo.GetBankAccount(userId)

	if err != nil {
		if err.Error() == "bank account not found" {
			return nil, 404, err
		} else {
			return nil, 500, err
		}
	}

	resp := []entity.BankAccountGetResponse{}

	for _, account := range res {
		tempAcc := entity.BankAccountGetResponse{
			ID:            account.ID,
			Name:          account.Name,
			AccountName:   account.AccountName,
			AccountNumber: account.AccountNumber,
		}

		resp = append(resp, tempAcc)
	}

	return resp, 200, nil
}
