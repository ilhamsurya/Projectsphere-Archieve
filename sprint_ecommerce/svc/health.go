package svc

import "github.com/ramadhan1445sprint/sprint_ecommerce/repo"

type HealthSvcInterface interface {
	GetStatus() (string, error)
}

func NewHealthSvc(repo repo.HealthRepoInterface) HealthSvcInterface {
	return &healthSvc{repo: repo}
}

type healthSvc struct {
	repo repo.HealthRepoInterface
}

func (s *healthSvc) GetStatus() (string, error) {
	status, err := s.repo.GetStatus()
	if err != nil {
		return "", err
	}

	return status, nil
}
