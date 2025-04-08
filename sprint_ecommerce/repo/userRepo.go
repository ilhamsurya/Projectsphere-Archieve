package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
)

type UserRepo interface {
	GetUser(username string) (*entity.User, error)
	CreateUser(name, username, password string) (string, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetUser(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, "SELECT id, name, username, password FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) CreateUser(name, username, password string) (string, error) {
	var id string
	row := r.db.QueryRow("INSERT INTO users (name, username, password) VALUES ($1, $2, $3) RETURNING id", name, username, password)

	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}
