package svc

import (
	"errors"
	"sync"

	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
	"github.com/ramadhan1445sprint/sprint_ecommerce/repo"
	"github.com/ramadhan1445sprint/sprint_ecommerce/utils"
)

type PaymentSvcInterface interface {
	CreatePayment(paymentReq *entity.PaymentCreateRequest, productId string) (int, error)
}

func NewPaymentSvc(repo repo.PaymentRepoInterface) PaymentSvcInterface {
	return &paymentSvc{repo: repo}
}

type paymentSvc struct {
	repo repo.PaymentRepoInterface
	mutex sync.Mutex
}

func (s *paymentSvc) CreatePayment(paymentReq *entity.PaymentCreateRequest, productId string) (int, error) {
	paymentReq.ProductID = &productId
	status, err := utils.ValidatePaymentRequest(paymentReq)

	if err != nil {
		return status, err
	}

	s.mutex.Lock()
    defer s.mutex.Unlock()

	productQty, err := s.repo.GetProductQty(productId)

	if err != nil {
		return 500, err
	}

	if productQty < *paymentReq.Quantity {
		return 400, errors.New("insufficient quantity")
	}

	payment := entity.Payment{
		ProductID: *paymentReq.ProductID,
		BankAccountID: *paymentReq.BankAccountID,
		PaymentProofImgUrl: *paymentReq.PaymentProofImgUrl,
		Quantity: *paymentReq.Quantity,
	}

	if err = s.repo.CreatePayment(&payment); err != nil {
		return 500, err
	}

	stock := productQty - payment.Quantity

	if err = s.repo.UpdateStock(productId, stock); err != nil {
		return 500 , err
	}

	return  200, nil
}
