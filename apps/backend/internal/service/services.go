package service

import (
	"github.com/srauf24/gardenjournal/internal/lib/job"
	"github.com/srauf24/gardenjournal/internal/repository"
	"github.com/srauf24/gardenjournal/internal/server"
)

type Services struct {
	Auth *AuthService
	Job  *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	return &Services{
		Job:  s.Job,
		Auth: authService,
	}, nil
}
