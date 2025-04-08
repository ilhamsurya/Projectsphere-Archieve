package svc

import (
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ramadhan1445sprint/sprint_ecommerce/config"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
	"github.com/ramadhan1445sprint/sprint_ecommerce/pkg/database"
	"github.com/ramadhan1445sprint/sprint_ecommerce/repo"
)

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func TestCreateBankAccount(t *testing.T) {
	config.LoadConfig("../.env")

	db, err := database.NewDatabase()
	if err != nil {
		t.Fatalf("failed to create a database connection: %v", err)
	}

	bankAccountRepo := repo.NewBankAccountRepo(db)
	bankAccountSvc := NewBankAccounthSvc(bankAccountRepo)

	testCases := []struct {
		name        		string
		input       		entity.BankAccountCreateRequest
		userID				string
		statusExpected 		int
		errExpected 		bool
	}{
		{"Test create bank success", entity.BankAccountCreateRequest{Name: strPtr("svc test"), AccountName: strPtr("svc taest name"), AccountNumber: strPtr("23123")}, "7d05089b-23a1-4c95-98f7-840f144428b3", 200, false},
		{"Test create bank failed", entity.BankAccountCreateRequest{Name: nil, AccountName: strPtr("svc taest name"), AccountNumber: strPtr("23123")}, "7d05089b-23a1-4c95-98f7-840f144428b3", 400, true},
		{"Test create bank failed", entity.BankAccountCreateRequest{Name: strPtr("svc test"), AccountName: nil, AccountNumber: strPtr("23123")}, "7d05089b-23a1-4c95-98f7-840f144428b3", 400, true},
		{"Test create bank failed", entity.BankAccountCreateRequest{Name: strPtr("svc test"), AccountName: strPtr("svc taest name"), AccountNumber: nil}, "7d05089b-23a1-4c95-98f7-840f144428b3", 400, true},
		{"Test create bank failed", entity.BankAccountCreateRequest{Name: strPtr(""), AccountName: strPtr("svc taest name"), AccountNumber: strPtr("23123")}, "7d05089b-23a1-4c95-98f7-840f144428b3", 400, true},
		{"Test create bank failed", entity.BankAccountCreateRequest{Name: strPtr("svc test"), AccountName: strPtr(""), AccountNumber: strPtr("23123")}, "7d05089b-23a1-4c95-98f7-840f144428b3", 400, true},
		{"Test create bank failed", entity.BankAccountCreateRequest{Name: strPtr("svc test"), AccountName: strPtr("dadawd"), AccountNumber: strPtr("0")}, "7d05089b-23a1-4c95-98f7-840f144428b3", 400, true},
		{"Test create bank failed", entity.BankAccountCreateRequest{Name: strPtr("svc test iadiajwdijaidjaiddawjiawjd"), AccountName: strPtr("dadawd"), AccountNumber: strPtr("0")}, "7d05089b-23a1-4c95-98f7-840f144428b3", 400, true},
		{"Test create bank failed", entity.BankAccountCreateRequest{Name: strPtr("svc test"), AccountName: strPtr("dkoakdoak oakwdoawkdo koakwdoakwdkkdoaw"), AccountNumber: strPtr("0")}, "7d05089b-23a1-4c95-98f7-840f144428b3", 400, true},
		{"Test create bank failed", entity.BankAccountCreateRequest{Name: strPtr("svc test"), AccountName: strPtr("dadawd"), AccountNumber: strPtr("123")}, "7d05089b-23a1-4c95-98f7-840f144428b3", 400, true},
		{"Test create bank failed", entity.BankAccountCreateRequest{Name: strPtr("svc test"), AccountName: strPtr("dadawd"), AccountNumber: strPtr("1234567890123456")}, "7d05089b-23a1-4c95-98f7-840f144428b3", 400, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, err := bankAccountSvc.CreateBankAccount(&tc.input, tc.userID)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error, but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}

				if status != 200 {
					t.Errorf("expected %d, but got %d", tc.statusExpected, status)
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

	bankAccountRepo := repo.NewBankAccountRepo(db)
	bankAccountSvc := NewBankAccounthSvc(bankAccountRepo)

	testCases := []struct {
		name        string
		input       string
		errExpected bool
	}{
		{"Test get bank account success", "7d05089b-23a1-4c95-98f7-840f144428b3", false},
		{"Test get bank account fail", "99999", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, status, err := bankAccountSvc.GetBankAccount(tc.input)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error, but no error")
				}
			} else {

				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}

				if status != 200 {
					t.Errorf("expected %d, but got %d", 200, status)
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

	bankAccountRepo := repo.NewBankAccountRepo(db)
	bankAccountSvc := NewBankAccounthSvc(bankAccountRepo)

	testCases := []struct {
		name        string
		input       entity.BankAccountUpdateRequest
		userId		string
		errExpected bool
	}{
		{"Test update bank account success", entity.BankAccountUpdateRequest{ID: strPtr("fd10ed3b-5d94-4c3f-97d2-deaf99644efd"), Name: strPtr("update svc 2"), AccountName: strPtr("dadang svc 2"), AccountNumber: strPtr("123454")}, "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", false},
		{"Test failed bank failed", entity.BankAccountUpdateRequest{ID: strPtr("fd10ed3b-5d94-4c3f-97d2-deaf99644efd"), Name: nil, AccountName: strPtr("svc taest name"), AccountNumber: strPtr("23123")}, "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", true},
		{"Test failed bank failed", entity.BankAccountUpdateRequest{ID: strPtr("fd10ed3b-5d94-4c3f-97d2-deaf99644efd"), Name: strPtr("svc test"), AccountName: nil, AccountNumber: strPtr("23123")}, "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", true},
		{"Test failed bank failed", entity.BankAccountUpdateRequest{ID: strPtr("fd10ed3b-5d94-4c3f-97d2-deaf99644efd"), Name: strPtr("svc test"), AccountName: strPtr("svc taest name"), AccountNumber: nil}, "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", true},
		{"Test failed bank failed", entity.BankAccountUpdateRequest{ID: strPtr("fd10ed3b-5d94-4c3f-97d2-deaf99644efd"), Name: strPtr(""), AccountName: strPtr("svc taest name"), AccountNumber: strPtr("23123")}, "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", true},
		{"Test failed bank failed", entity.BankAccountUpdateRequest{ID: strPtr("fd10ed3b-5d94-4c3f-97d2-deaf99644efd"), Name: strPtr("svc test"), AccountName: strPtr(""), AccountNumber: strPtr("23123")}, "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", true},
		{"Test update bank failed", entity.BankAccountUpdateRequest{ID: strPtr("fd10ed3b-5d94-4c3f-97d2-deaf99644efd"), Name: strPtr("svc test"), AccountName: strPtr("dadawd"), AccountNumber: strPtr("0")}, "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", true},
		{"Test update bank failed", entity.BankAccountUpdateRequest{ID: strPtr("fd10ed3b-5d94-4c3f-97d2-deaf99644efd"), Name: strPtr("svc test iadiajwdijaidjaiddawjiawjd"), AccountName: strPtr("dadawd"), AccountNumber: strPtr("0")}, "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", true},
		{"Test update bank failed", entity.BankAccountUpdateRequest{ID: strPtr("fd10ed3b-5d94-4c3f-97d2-deaf99644efd"), Name: strPtr("svc test"), AccountName: strPtr("dkoakdoak oakwdoawkdo koakwdoakwdkkdoaw"), AccountNumber: strPtr("0")}, "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", true},
		{"Test update bank failed", entity.BankAccountUpdateRequest{ID: strPtr("fd10ed3b-5d94-4c3f-97d2-deaf99644efd"), Name: strPtr("svc test"), AccountName: strPtr("dadawd"), AccountNumber: strPtr("123")}, "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", true},
		{"Test update bank failed", entity.BankAccountUpdateRequest{ID: strPtr("fd10ed3b-5d94-4c3f-97d2-deaf99644efd"), Name: strPtr("svc test"), AccountName: strPtr("dadawd"), AccountNumber: strPtr("1234567890123456")}, "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", true},
		{"Test update bank failed", entity.BankAccountUpdateRequest{ID: strPtr("dwdwd"), Name: strPtr("svc test"), AccountName: strPtr("dadawd"), AccountNumber: strPtr("1234567890123456")}, "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, err := bankAccountSvc.UpdateBankAccount(&tc.input, tc.userId)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error, but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}

				if status != 200 {
					t.Errorf("expected %d, but got %d", 200, status)
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

	bankAccountRepo := repo.NewBankAccountRepo(db)
	bankAccountSvc := NewBankAccounthSvc(bankAccountRepo)

	testCases := []struct {
		name        string
		input       string
		userId      string
		errExpected bool
	}{
		{"Test delete bank account success", "fd10ed3b-5d94-4c3f-97d2-deaf99644efd", "e2bb8ffa-d01c-41d0-a438-152f6469b4f2", false},
		{"Test delete bank account failed", "99999", "e2bb8ffa-d01c-41d0-a438-152f6469b4f1", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, err := bankAccountSvc.DeleteBankAccount(tc.input, tc.userId)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error, but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}

				if status != 200 {
					t.Errorf("expected %d, but got %d", 200, status)
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

	paymentRepo := repo.NewPaymentRepo(db)
	paymentSvc := NewPaymentSvc(paymentRepo)

	testCases := []struct {
		name        string
		input       entity.PaymentCreateRequest
		productId   string
		errExpected bool
	}{
		{"Test create payment success", entity.PaymentCreateRequest{BankAccountID: strPtr("02086ff6-7df3-44a8-aebb-4906f0360c39"), PaymentProofImgUrl: strPtr("https://example.com/image.jpg"), Quantity: intPtr(10)}, "43826207-2a72-40c5-a696-35cc66c32e2e", false},
		{"Test create payment failed 1", entity.PaymentCreateRequest{BankAccountID: nil, PaymentProofImgUrl: strPtr("https://example.com/image.jpg"), Quantity: intPtr(10)}, "43826207-2a72-40c5-a696-35cc66c32e2e", true},
		{"Test create payment failed 2", entity.PaymentCreateRequest{BankAccountID: strPtr("02086ff6-7df3-44a8-aebb-4906f0360c39"), PaymentProofImgUrl: strPtr("https://example.com/image.jpg"), Quantity: intPtr(10)}, "43826207-2a72-40c5-a696-35cc66c32e2e", true},
		{"Test create payment failed 3", entity.PaymentCreateRequest{BankAccountID: strPtr("02086ff6-7df3-44a8-aebb-4906f0360c39"), PaymentProofImgUrl: nil, Quantity: intPtr(10)}, "43826207-2a72-40c5-a696-35cc66c32e2e", true},
		{"Test create payment failed 4", entity.PaymentCreateRequest{BankAccountID: strPtr("02086ff6-7df3-44a8-aebb-4906f0360c39"), PaymentProofImgUrl: strPtr("https://example.com/image.jpg"), Quantity: nil}, "43826207-2a72-40c5-a696-35cc66c32e2e", true},
		{"Test create payment failed 5", entity.PaymentCreateRequest{BankAccountID: strPtr("02086ff6-7df3-44a8-aebb-4906f0360c39"), PaymentProofImgUrl: strPtr("https://example.com/image.jpg"), Quantity: intPtr(0)}, "43826207-2a72-40c5-a696-35cc66c32e2e", true},
		{"Test create payment failed 6", entity.PaymentCreateRequest{BankAccountID: strPtr("02086ff6-7df3-44a8-aebb-4906f0360c39"), PaymentProofImgUrl: strPtr("owkeokweo"), Quantity: intPtr(0)}, "43826207-2a72-40c5-a696-35cc66c32e2e", true},
		{"Test create payment failed 7", entity.PaymentCreateRequest{BankAccountID: strPtr("43826207-2a72-40c5-a696-35cc66c32e2e"), PaymentProofImgUrl: strPtr("https://example.com/image.jpg"), Quantity: intPtr(10)}, "43826207-2a72-40c5-a696-35cc66c32e2e", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, err := paymentSvc.CreatePayment(&tc.input, tc.productId)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}

				if status != 200 {
					t.Errorf("expected %d, but got %d", 200, status)
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

	stockRepo := repo.NewStockRepo(db)
	stockSvc := NewStockSvc(stockRepo)

	testCases := []struct {
		name        string
		input       entity.StockUpdateRequest
		productId   string
		userId 		string
		errExpected bool
	}{
		{"Test update stock success", entity.StockUpdateRequest{Stock: intPtr(111)}, "43826207-2a72-40c5-a696-35cc66c32e2e", "6ca7abef-4f8a-4613-9a3f-77e9b34c8c49",false},
		{"Test update stock failed", entity.StockUpdateRequest{Stock: intPtr(111)}, "", "6ca7abef-4f8a-4613-9a3f-77e9b34c8c49",true},
		{"Test update stock failed 1", entity.StockUpdateRequest{Stock: intPtr(111)}, "", "6ca7abef-4f8a-4613-9a3f-77e9b34c8c49",true},
		{"Test update stock failed 2", entity.StockUpdateRequest{Stock: nil}, "43826207-2a72-40c5-a696-35cc66c32e2e", "6ca7abef-4f8a-4613-9a3f-77e9b34c8c49",true},
		{"Test update stock failed 3", entity.StockUpdateRequest{Stock: intPtr(-1)},"43826207-2a72-40c5-a696-35cc66c32e2e", "6ca7abef-4f8a-4613-9a3f-77e9b34c8c49",true},
		
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, err := stockSvc.UpdateStock(&tc.input, tc.productId, tc.userId)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}

				if status != 200 {
					t.Errorf("expected %d, but got %d", 200, status)
				}
			}
		})
	}
}
