package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
)

type PaymentRepoInterface interface {
	CreatePayment(payment *entity.Payment) error
	GetProductQty(productId string) (int, error)
	UpdateStock(productId string, stock int) error
}

func NewPaymentRepo(db *sqlx.DB) PaymentRepoInterface {
	return &paymentRepo{db: db}
}

type paymentRepo struct {
	db *sqlx.DB
}

func (r *paymentRepo) CreatePayment(payment *entity.Payment) error {
	fmt.Printf("%#v", payment)
	_, err := r.db.Exec("INSERT INTO payments (product_id, bank_account_id, quantity, payment_proof_image_url) values ($1, $2, $3, $4)", payment.ProductID, payment.BankAccountID, payment.Quantity, payment.PaymentProofImgUrl)

	if err != nil {
		return err
	}

	return nil
}

func (r *paymentRepo) GetProductQty(productId string) (int, error) {

	var res *int

	err := r.db.Get(&res, "select stock from products where id = $1", productId)

	if err != nil {
		return 0, err
	}

	return *res, nil
}

func (r *paymentRepo) UpdateStock(productId string, stock int) error {
	_, err := r.db.Exec("update products set stock = $1 where id = $2", stock, productId)

	if err != nil {
		return err
	}

	return nil
}
