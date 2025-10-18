package repository

import "github.com/srauf24/gardenjournal/internal/server"

type PlantRepository struct {
    server *server.Server
}

func NewPlantRepository(s *server.Server) *PlantRepository {
    return &PlantRepository{server: server}
}