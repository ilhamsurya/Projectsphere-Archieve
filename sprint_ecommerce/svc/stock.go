package svc

import (

	"github.com/google/uuid"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
	"github.com/ramadhan1445sprint/sprint_ecommerce/repo"
	"github.com/ramadhan1445sprint/sprint_ecommerce/utils"
)

type StockSvcInterface interface {
	UpdateStock(productReq *entity.StockUpdateRequest, productId string, userId string) (int, error)
}

func NewStockSvc(repo repo.StockRepoInterface) StockSvcInterface {
	return &stockSvc{repo: repo}
}

type stockSvc struct {
	repo repo.StockRepoInterface
}

func (s *stockSvc) UpdateStock(productReq *entity.StockUpdateRequest, productId string, userId string) (int, error) {
	_, err := uuid.Parse(productId)

	if err != nil {
		return 500, err
	}
	productReq.ProductID = &productId

	status, err := utils.ValidateStockUpdateRequest(productReq)

	if err != nil {
		return status, err
	}

	if err := s.repo.GetProductById(productId); err != nil {
		if err.Error() == "product not found" {
			return 404, err
		} else {
			return 500, err
		}
	}

	if err := s.repo.CheckProductByUserId(userId, productId); err != nil {
		if err.Error() == "not allowed" {
			return 403, err
		} else {
			return 500, err
		}
	}

	product := entity.Stock{
		ID:    *productReq.ProductID,
		Stock: *productReq.Stock,
	}

	err = s.repo.UpdateStock(&product)
	if err != nil {
		return 500, err
	}

	return 200, nil
}
