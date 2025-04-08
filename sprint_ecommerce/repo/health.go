package repo

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type HealthRepoInterface interface {
	GetStatus() (string, error)
}

func NewHealthRepo(db *sqlx.DB) HealthRepoInterface {
	return &healthRepo{db: db}
}

type healthRepo struct {
	db *sqlx.DB
}

func (r *healthRepo) GetStatus() (string, error) {
	var res int
	r.db.Get(&res, "SELECT 1")

	if res == 1 {
		return "ok", nil
	}

	return "", errors.New("something went wrong")
}
