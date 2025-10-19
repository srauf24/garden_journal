package repository

import "github.com/srauf24/gardenjournal/internal/server"

type PlantRepository struct {
    server *server.Server
}

func NewPlantRepository(server *server.Server) *PlantRepository {
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

func (r *PlantRepository) GetPlantByID(ctx context.Context, userID string, plantID uuid.UUID) (*plant.Plant, error) {
	stmt := `
	SELECT
		p.*,

		-- Aggregate all observations for this plant
		COALESCE(
			jsonb_agg(
				to_jsonb(camel(obs))
				ORDER BY obs.date DESC
			) FILTER (
				WHERE obs.id IS NOT NULL
			),
			'[]'::JSONB
		) AS observations

	FROM
		plants p
		LEFT JOIN observations obs
			ON obs.plant_id = p.id
			AND obs.user_id = @user_id

	WHERE
		p.id = @plant_id
		AND p.user_id = @user_id

	GROUP BY
		p.id
	HAVING
		p.id IS NOT NULL
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":  userID,
		"plant_id": plantID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get plant by id query for user_id=%s plant_id=%s: %w", userID, plantID, err)
	}

	plantItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[plant.Plant])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:plants for user_id=%s plant_id=%s: %w", userID, plantID, err)
	}

	return &plantItem, nil
}
