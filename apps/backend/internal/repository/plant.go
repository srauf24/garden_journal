package repository

import "github.com/srauf24/gardenjournal/internal/server"

type PlantRepository struct {
    server *server.Server
}

func NewPlantRepository(s *server.Server) *PlantRepository {
    return &PlantRepository{server: server}
}
func (r * PlantRepository) CreatePlant(ctx context.Context, userID string, payload *plant.CreatePlantPayload) (*plant.Plant, error) {
stmt := `
		INSERT INTO
			plants (
				user_id,
				name,
				species,
				location,
				planted_date,
				notes,
				metadata
			)
		VALUES
			(
				@user_id,
				@name,
				@species,
				@location,
				@planted_date,
				@notes,
				@metadata
			)
		RETURNING
		*
	`
	// use server.db to execute the query
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":        userID,
		"name":           payload.Name,
		"species":        payload.Species,
		"location":       payload.Location,
		"planted_date":   payload.PlantedDate,
		"notes":          payload.Notes,
		"metadata":       payload.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create plant query for user_id=%s title=%s: %w", userID, payload.Title, err)
	}
// use pgx library to deserialize row into a struct from the data base (collect one row)
	plantItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[plant.Plant])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:plants for user_id=%s title=%s: %w", userID, payload.Title, err)
	}

	return &plantItem, nil
}