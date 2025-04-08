package svc

import (
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
	"github.com/ramadhan1445sprint/sprint_ecommerce/repo"
	"github.com/ramadhan1445sprint/sprint_ecommerce/utils"
)

type BankAccountSvcInterface interface {
	GetBankAccount(userId string) ([]entity.BankAccountGetResponse, int, error)
	CreateBankAccount(bankAccount *entity.BankAccountCreateRequest, userId string) (int, error)
	UpdateBankAccount(bankAccount *entity.BankAccountUpdateRequest, userId string) (int, error)
	DeleteBankAccount(bankAccountID string, userId string) (int, error)
}

func NewBankAccounthSvc(repo repo.BankAccountRepoInterface) BankAccountSvcInterface {
	return &bankAccountSvc{repo: repo}
}

type bankAccountSvc struct {
	repo repo.BankAccountRepoInterface
}

func (s *bankAccountSvc) GetBankAccount(userId string) ([]entity.BankAccountGetResponse, int, error) {
	res, err := s.repo.GetBankAccount(userId)

	if err != nil {
		if err.Error() == "bank account not found" {
			return nil, 404, err
		}else {
			return nil, 500, err		
		}
	}

	resp := []entity.BankAccountGetResponse{}

	for _, account := range res {
		tempAcc := entity.BankAccountGetResponse{
			ID: account.ID,
			Name: account.Name,
			AccountName: account.AccountName,
			AccountNumber: account.AccountNumber,
		}

		resp = append(resp, tempAcc)
	} 

	return resp, 200, nil
}

func (s *bankAccountSvc) CreateBankAccount(req *entity.BankAccountCreateRequest, userId string) (int, error) {
	status, err := utils.ValidateCreateBankRequest(req)

	if err != nil {
		return status, err
	}

	bankAccount := entity.BankAccount{
		UserID: userId,
		Name: *req.Name,
		AccountName: *req.AccountName,
		AccountNumber: *req.AccountNumber,
	}

	if err = s.repo.CreateBankAccount(&bankAccount); err != nil {
		return 500, err
	}

	return 200, nil
}

func (s *bankAccountSvc) UpdateBankAccount(req *entity.BankAccountUpdateRequest, userId string) (int, error) {
	status, err := utils.ValidateUpdateBankRequest(req)

	if err != nil {
		return status, err
	}

	bankAccount := entity.BankAccount{
		ID: *req.ID,
		Name: *req.Name,
		AccountName: *req.AccountName,
		AccountNumber: *req.AccountNumber,
	}

	if err := s.repo.CheckBankAccountByUserId(userId, *req.ID); err != nil {
		if err.Error() == "not allowed" {
			return 403, err
		} else {
			return 500, err
		}
	}

	if err = s.repo.UpdateBankAccount(&bankAccount); err != nil {
		return 500, err
	}

	return 200, nil
}

func (s *bankAccountSvc) DeleteBankAccount(bankAccountID string, userId string) (int, error) {

	if err := s.repo.CheckBankAccountByUserId(userId, bankAccountID); err != nil {
		if err.Error() == "not allowed" {
			return 403, err
		} else {
			return 500, err
		}
	}

	if err := s.repo.DeleteBankAccount(bankAccountID); err != nil {
		var status int
		if err.Error() == "bank account not found" {
			status = 400
		}else {
			status = 500
		}
		return status, err
	}

	return 200, nil
}
