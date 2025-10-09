package repository

import "github.com/srauf24/gardenjournal/internal/server"

type Repositories struct{}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}
