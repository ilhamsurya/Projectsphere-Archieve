package svc

import (
	"github.com/ramadhan1445sprint/sprint_ecommerce/crypto"
	"github.com/ramadhan1445sprint/sprint_ecommerce/customErr"
	"github.com/ramadhan1445sprint/sprint_ecommerce/entity"
	"github.com/ramadhan1445sprint/sprint_ecommerce/repo"
)

type UserSvc interface {
	RegisterUser(*entity.User) (string, error)
	Login(entity.Credential) (*entity.User, string, error)
}

type userSvc struct {
	repo repo.UserRepo
}

func NewUserSvc(repo repo.UserRepo) UserSvc {
	return &userSvc{repo}
}

func (s *userSvc) RegisterUser(user *entity.User) (string, error) {
	if user.Name == "" {
		return "", customErr.NewBadRequestError("name is required")
	}

	if len(user.Name) < 5 || len(user.Name) > 50 {
		return "", customErr.NewBadRequestError("name must be between 5 and 50 characters")
	}

	if user.Username == "" {
		return "", customErr.NewBadRequestError("username is required")
	}

	if len(user.Username) < 5 || len(user.Username) > 15 {
		return "", customErr.NewBadRequestError("username must be between 5 and 15 characters")
	}

	if user.Password == "" {
		return "", customErr.NewBadRequestError("password is required")
	}

	if len(user.Password) < 5 || len(user.Password) > 15 {
		return "", customErr.NewBadRequestError("password must be between 5 and 15 characters")
	}

	existingUser, err := s.repo.GetUser(user.Username)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return "", err
		}
	}
	if existingUser != nil {
		return "", customErr.NewConflictError("username already exists")
	}

	// bcrypt user password
	hashedPassword, err := crypto.GenerateHashedPassword(user.Password)
	if err != nil {
		return "", err
	}

	uid, err := s.repo.CreateUser(user.Name, user.Username, hashedPassword)
	if err != nil {
		return "", err
	}

	token, err := crypto.GenerateToken(uid, user.Username, user.Name)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userSvc) Login(creds entity.Credential) (*entity.User, string, error) {
	if creds.Username == "" {
		return nil, "", customErr.NewBadRequestError("username is required")
	}

	if len(creds.Username) < 5 || len(creds.Username) > 15 {
		return nil, "", customErr.NewBadRequestError("username must be between 5 and 15 characters")
	}

	if creds.Password == "" {
		return nil, "", customErr.NewBadRequestError("password is required")
	}

	if len(creds.Password) < 5 || len(creds.Password) > 15 {
		return nil, "", customErr.NewBadRequestError("password must be between 5 and 15 characters")
	}

	user, err := s.repo.GetUser(creds.Username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, "", customErr.NewNotFoundError("username not found")
		}
		return nil, "", err
	}

	err = crypto.VerifyPassword(creds.Password, user.Password)
	if err != nil {
		return nil, "", customErr.NewBadRequestError("wrong password!")
	}

	token, err := crypto.GenerateToken(user.Id, user.Username, user.Name)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
