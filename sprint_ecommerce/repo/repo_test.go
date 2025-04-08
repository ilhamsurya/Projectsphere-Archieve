package repo

import (
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ramadhan1445sprint/sprint_ecommerce/config"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
	"github.com/ramadhan1445sprint/sprint_ecommerce/pkg/database"
)

func TestCreateBankAccount(t *testing.T) {
	config.LoadConfig("../.env")

	db, err := database.NewDatabase()
	if err != nil {
		t.Fatalf("failed to create a database connection: %v", err)
	}

	bankAccountRepo := NewBankAccountRepo(db)

	testCases := []struct {
		name        string
		input       entity.BankAccount
		errExpected bool
	}{
		{"Test create bank", entity.BankAccount{UserID: "7d05089b-23a1-4c95-98f7-840f144428b3", Name: "BCA", AccountName: "Ilham Nuryanto", AccountNumber: "1234567"}, false},
		{"Test create bank 2", entity.BankAccount{UserID: "7d05089b-23a1-4c95-98f7-840f144428b3", Name: "Mandiri", AccountName: "Ilham Nuryanto", AccountNumber: "432145"}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := bankAccountRepo.CreateBankAccount(&tc.input)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error, but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}

}
func TestGetBankAccount(t *testing.T) {
	config.LoadConfig("../.env")

	db, err := database.NewDatabase()
	if err != nil {
		t.Fatalf("failed to create a database connection: %v", err)
	}

	bankAccountRepo := NewBankAccountRepo(db)

	testCases := []struct {
		name        string
		input       string
		errExpected bool
	}{
		{"Test get bank account found", "7739f6f4-8e6e-42b4-bed0-d87ca4499353", false},
		{"Test get bank account not found", "99999", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := bankAccountRepo.GetBankAccount(tc.input)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error, but no error")
				}
			} else {
				if len(res) == 0 {
					t.Errorf("Expected exist, but got: %d", len(res))
				}

				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

func TestUpdateBankAccount(t *testing.T) {
	config.LoadConfig("../.env")

	db, err := database.NewDatabase()
	if err != nil {
		t.Fatalf("failed to create a database connection: %v", err)
	}

	bankAccountRepo := NewBankAccountRepo(db)

	testCases := []struct {
		name        string
		input       entity.BankAccount
		errExpected bool
	}{
		{"Test update bank account success", entity.BankAccount{ID: "02086ff6-7df3-44a8-aebb-4906f0360c39", Name: "BCA updated", AccountName: "ilham updated", AccountNumber: "123123"}, false},
		{"Test update bank account failed", entity.BankAccount{ID: "99999", Name: "BCA", AccountName: "dadang", AccountNumber: "1234"}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := bankAccountRepo.UpdateBankAccount(&tc.input)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error, but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}
func TestDeleteBankAccount(t *testing.T) {
	config.LoadConfig("../.env")

	db, err := database.NewDatabase()
	if err != nil {
		t.Fatalf("failed to create a database connection: %v", err)
	}

	bankAccountRepo := NewBankAccountRepo(db)

	testCases := []struct {
		name        string
		input       string
		errExpected bool
	}{
		{"Test delete bank account success", "b89a6209-ba1d-4b82-95fe-423f1a43653d", false},
		{"Test delete bank account failed", "9999", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := bankAccountRepo.DeleteBankAccount(tc.input)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error, but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}
func TestCreatePayment(t *testing.T) {
	config.LoadConfig("../.env")

	db, err := database.NewDatabase()
	if err != nil {
		t.Fatalf("failed to create a database connection: %v", err)
	}

	paymentRepo := NewPaymentRepo(db)

	testCases := []struct {
		name        string
		input       entity.Payment
		errExpected bool
	}{
		{"Test create payment success", entity.Payment{ProductID: "43826207-2a72-40c5-a696-35cc66c32e2e", BankAccountID: "02086ff6-7df3-44a8-aebb-4906f0360c39", PaymentProofImgUrl: "url", Quantity: 10}, false},
		{"Test create payment failed", entity.Payment{ProductID: "dw2", BankAccountID: "02086ff6-7df3-44a8-aebb-4906f0360c39", PaymentProofImgUrl: "url", Quantity: 10}, true},
		{"Test create payment failed", entity.Payment{ProductID: "43826207-2a72-40c5-a696-35cc66c32e2e", BankAccountID: "ww", PaymentProofImgUrl: "url", Quantity: 10}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := paymentRepo.CreatePayment(&tc.input)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

func TestUpdateStock(t *testing.T) {
	config.LoadConfig("../.env")

	db, err := database.NewDatabase()
	if err != nil {
		t.Fatalf("failed to create a database connection: %v", err)
	}

	stockRepo := NewStockRepo(db)

	testCases := []struct {
		name        string
		input       entity.Stock
		errExpected bool
	}{
		{"Test create payment success", entity.Stock{ID: "36b11d41-144e-478e-876f-e6118e4b23db", Stock: 999}, false},
		{"Test create payment failed", entity.Stock{ID: "dwad", Stock: 2}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := stockRepo.UpdateStock(&tc.input)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

func TestCheckProductByUserId(t *testing.T) {
	config.LoadConfig("../.env")

	db, err := database.NewDatabase()
	if err != nil {
		t.Fatalf("failed to create a database connection: %v", err)
	}

	stockRepo := NewStockRepo(db)

	testCases := []struct {
		name        string
		userId      string
		productId 	string
		errExpected bool
	}{
		{"Test create payment success", "6ca7abef-4f8a-4613-9a3f-77e9b34c8c49", "97ed68e2-221e-4e0f-885f-d37ef90919c2", false},
		{"Test create payment failed", "dwa", "97ed68e2-221e-4e0f-885f-d37ef90919c2", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := stockRepo.CheckProductByUserId(tc.userId, tc.productId)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}
