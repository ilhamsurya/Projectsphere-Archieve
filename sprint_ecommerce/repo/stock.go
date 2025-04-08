package repo

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
)

type StockRepoInterface interface {
	UpdateStock(product *entity.Stock) error
	CheckProductByUserId(userId string, productId string) error
	GetProductById(productId string) error
}

func NewStockRepo(db *sqlx.DB) StockRepoInterface {
	return &stockRepo{db: db}
}

type stockRepo struct {
	db *sqlx.DB
}

func (r *stockRepo) UpdateStock(product *entity.Stock) error {
	res, err := r.db.Exec("UPDATE products SET stock = $1 where id = $2", product.Stock, product.ID)

	if err != nil {
		return err
	}

	rowsEffected, _ := res.RowsAffected()

	if rowsEffected == 0 {
		return errors.New("bank account not found")
	}

	return nil
}

func (r *stockRepo) CheckProductByUserId(userId string, productId string) error {
	var res *int

	err := r.db.Get(&res, "select count(id) from products where user_id = $1 and id = $2", userId, productId)

	if err != nil {
		return err
	}

	if *res == 0 {
		return errors.New("not allowed")
	}

	return nil
}

func (r *stockRepo) GetProductById(productId string) error {
	var res *int

	err := r.db.Get(&res, "select count(id) from products where id = $1", productId)

	if err != nil {
		return err
	}

	if *res == 0 {
		return errors.New("product not found")
	}

	return nil
}
