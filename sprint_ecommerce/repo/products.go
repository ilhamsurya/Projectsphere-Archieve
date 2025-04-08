package repo

// import (

// 	"github.com/jmoiron/sqlx"
// 	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
// )

// type PaymentRepoInterface interface {
// 	CreatePayment(payment *entity.Payment) error
// }

// func NewPaymentRepo(db *sqlx.DB) PaymentRepoInterface {
// 	return &paymentRepo{db: db}
// }

// type paymentRepo struct {
// 	db *sqlx.DB
// }

// func (r *paymentRepo) CreatePayment(payment *entity.Payment) error {
// 	_, err := r.db.Exec("INSERT INTO payments (product_id, bank_account_id, quantity, payment_proof_image_url) values ($1, $2, $3, $4)", payment.ProductID, payment.BankAccountID, payment.Quantity, payment.PaymentProofImgUrl)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
