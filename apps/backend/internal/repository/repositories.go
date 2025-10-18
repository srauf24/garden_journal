package repository

import "github.com/srauf24/gardenjournal/internal/server"

type Repositories struct{
    Plant * PlantRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
        Plant: NewPlantRepository(s),
	}
}

// “Create a new Repositories struct with the Plant field initialized,
// and return a pointer to that struct.”
// Returning a pointer (*Repositories) is idiomatic in Go when:
//
// You want to modify the struct later (e.g. add repositories dynamically).
//
// You want to avoid copying a large struct.
//
// You want methods with pointer receivers (func (r *Repositories) DoSomething()) to work on the same instance.