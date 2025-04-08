package repo

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
)

type RepoInterface interface {
	CreateProduct(product entity.Product) error
	GetDetailProduct(id uuid.UUID) (entity.Product, error)
	UpdateProduct(product entity.Product) error
	DeleteProduct(id uuid.UUID, userId uuid.UUID) error
	GetListProduct(keys entity.Key, userId uuid.UUID) ([]entity.Product, error)
	GetProductSoldTotal(userId uuid.UUID) (entity.ProductPayment, error)
	GetBankAccount(userId string) ([]entity.BankAccount, error)
}

func NewRepo(db *sqlx.DB) RepoInterface {
	return &repo{db: db}
}

type repo struct {
	db *sqlx.DB
}

func (r *repo) CreateProduct(product entity.Product) error {
	query := `INSERT INTO products (user_id, name, price, stock, image_url, condition, is_purchasable, tags)
               VALUES (:user_id, :name, :price, :stock, :image_url, :condition, :is_purchasable, :tags)`

	_, err := r.db.NamedExec(query, &product)
	if err != nil {
		log.Println("Error executing query:", err)
		return err
	}

	return nil
}

func (r *repo) GetDetailProduct(id uuid.UUID) (entity.Product, error) {
	var product entity.Product
	var tags pgtype.VarcharArray
	query := "SELECT tags FROM products WHERE id = $1"

	// Query for a single row
	err := r.db.QueryRowx(query, id).Scan(&tags)
	if err != nil {
		log.Println("Error executing query:", err)
		return product, err
	}

	// Extract tags from VarcharArray
	var tagsSlice []string
	for _, tag := range tags.Elements {
		if tag.Status != pgtype.Null {
			tagsSlice = append(tagsSlice, string(tag.String))
		}
	}

	query1 := "SELECT id, user_id, name, price, stock, image_url, condition, is_purchasable FROM products WHERE id = $1"
	err = r.db.QueryRowx(query1, id).Scan(
		&product.ID,
		&product.UserID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.ImageUrl,
		&product.Condition,
		&product.IsPurchasable,
	)
	if err != nil {
		log.Println("Error executing query:", err)
		return product, err
	}

	total, err := GetPurchaseCount(product.ID, r)
	if err != nil {
		log.Println("Error executing query:", err)
		return product, err
	}

	product.Tags = tagsSlice
	product.PurchaseCount = total

	return product, nil
}

func GetPurchaseCount(id uuid.UUID, r *repo) (int, error) {
	var total int
	query := fmt.Sprintf(`SELECT COALESCE(SUM(quantity), 0) FROM payments WHERE product_id = '%s'`, id)

	// Query sum quantity
	err := r.db.QueryRow(query).Scan(&total)
	if err != nil {
		log.Println("Error executing query:", err)
		return 0, err
	}

	return total, err
}

func (r *repo) GetProductSoldTotal(userId uuid.UUID) (entity.ProductPayment, error) {
	var productPayment entity.ProductPayment

	query := fmt.Sprintf(`select u."name", COALESCE(SUM(quantity), 0) as totalSold from products p
											inner join users u on p.user_id = u.id
											left join payments p2 ON p.id = p2.product_id
											where u.id = '%s'
											group by u.id`, userId)

	// Query sum quantity
	err := r.db.QueryRow(query).Scan(&productPayment.Name, &productPayment.TotalSold)
	if err != nil {
		log.Println("Error executing query:", err)
		return productPayment, err
	}

	return productPayment, err
}

func GetProductUser(userId uuid.UUID, id uuid.UUID, r *repo) error {
	var idProduct string
	query := fmt.Sprintf(`SELECT id FROM products WHERE user_id = '%s' and id = '%s'`, userId, id)

	err := r.db.QueryRow(query).Scan(&idProduct)
	if err != nil {
		log.Println("Error executing query:", err)
		return errors.New("product is not owned by this user")
	}

	return nil
}

func (r *repo) UpdateProduct(product entity.Product) error {
	err := GetProductUser(product.UserID, product.ID, r)
	if err != nil {
		return err
	}

	query := `UPDATE products set name = :name, price = :price, stock = :stock, image_url = :image_url,
							condition = :condition, is_purchasable = :is_purchasable, tags = :tags
							WHERE id = :id`

	res, err := r.db.NamedExec(query, &product)

	if err != nil {
		log.Println("Error executing query:", err)
		return err
	}

	rowsEffected, _ := res.RowsAffected()

	if rowsEffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *repo) DeleteProduct(id uuid.UUID, userId uuid.UUID) error {
	err := GetProductUser(userId, id, r)
	if err != nil {
		return err
	}
	query := `DELETE FROM products WHERE id = $1`

	res, err := r.db.Exec(query, id)
	if err != nil {
		log.Println("Error executing query:", err)
		return err
	}

	rowsEffected, _ := res.RowsAffected()

	if rowsEffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *repo) GetListProduct(keys entity.Key, userId uuid.UUID) ([]entity.Product, error) {
	// var tags pgtype.VarcharArray
	var conditions []string

	if keys.MaxPrice != nil && keys.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price BETWEEN %.2f AND %.2f", *keys.MinPrice, *keys.MaxPrice))
	}

	if keys.Condition != nil {
		conditions = append(conditions, fmt.Sprintf("condition = '%s'", *keys.Condition))
	}

	if keys.ShowEmptyStock != nil {
		if *keys.ShowEmptyStock {
			conditions = append(conditions, "stock = 0")
		} else {
			conditions = append(conditions, "stock > 0")
		}
	}

	if len(keys.Tags) > 0 {
		var tagConditions []string
		for _, tag := range keys.Tags {
			tagConditions = append(tagConditions, fmt.Sprintf("'%s'", tag))
		}
		conditions = append(conditions, fmt.Sprintf("ARRAY[%s] && tags", strings.Join(tagConditions, ",")))
	}

	if keys.Search != nil {
		conditions = append(conditions, fmt.Sprintf("name LIKE '%%%s%%'", *keys.Search))
	}

	if keys.UserOnly != nil {
		if *keys.UserOnly {
			conditions = append(conditions, fmt.Sprintf("user_id = '%s'", userId))
		}
	}

	// Check if any conditions were provided
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	orderByClause := ""
	if keys.OrderBy != nil && keys.SortBy != nil {
		orderByClause = fmt.Sprintf("ORDER BY %s %s", *keys.SortBy, *keys.OrderBy)
	}

	limitClause := ""
	if keys.Limit != nil {
		limitClause = fmt.Sprintf("LIMIT %d", *keys.Limit)
	}

	offsetClause := ""
	if keys.Offset != nil {
		offsetClause = fmt.Sprintf("OFFSET %d", *keys.Offset)
	}

	query := fmt.Sprintf(`
        SELECT id, user_id, name, price, stock, image_url, condition, is_purchasable
        FROM products
				%s
        %s
        %s
        %s`, whereClause, orderByClause, limitClause, offsetClause)

	// Execute the query
	var products []entity.Product
	err := r.db.Select(&products, query)
	if err != nil {
		log.Println("Error executing query:", err)
		return products, err
	}

	var productIDs []string
	for _, p := range products {
		productIDs = append(productIDs, fmt.Sprintf("'%s'", p.ID))
	}
	inClause := strings.Join(productIDs, ", ")

	if len(products) > 0 {
		var tags []pgtype.VarcharArray
		query1 := fmt.Sprintf(`SELECT tags FROM products WHERE id in (%s)`, inClause)

		// Query for a single row
		err = r.db.Select(&tags, query1)
		if err != nil {
			log.Println("Error executing query:", err)
			return products, err
		}

		// Extract tags from VarcharArray
		for i, tagArray := range tags {
			var tagsSlice []string
			for _, tag := range tagArray.Elements {
				if tag.Status != pgtype.Null {
					tagsSlice = append(tagsSlice, string(tag.String))
					products[i].Tags = tagsSlice
				}
			}
		}
	}

	return products, nil
}

func (r *repo) GetBankAccount(userId string) ([]entity.BankAccount, error) {
	var res []entity.BankAccount

	r.db.Select(&res, "SELECT id, bank_name, account_name, account_number from bank_accounts where user_id = $1", userId)

	return res, nil
}
