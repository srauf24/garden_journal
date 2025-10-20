package repository

import "github.com/srauf24/gardenjournal/internal/server"
import "github.com/srauf24/gardenjournal/internal/model/plant"

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

func (r *PlantRepository) GetPlantByID(ctx context.Context,userID string,plantID uuid.UUID,) (*plant.PopulatedPlant, error) {
	stmt := `
	SELECT
		p.*,

		-- Aggregate all observations for this plant
		COALESCE(
			jsonb_agg(
				to_jsonb(camel(obs))
				ORDER BY obs.date DESC
			) FILTER (WHERE obs.id IS NOT NULL),
			'[]'::JSONB
		) AS observations
	FROM plants p
	LEFT JOIN observations obs
		ON obs.plant_id = p.id
		AND obs.user_id = @user_id
	WHERE
		p.id = @plant_id
		AND p.user_id = @user_id
	GROUP BY p.id
	HAVING p.id IS NOT NULL
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":  userID,
		"plant_id": plantID,
	})
	if err != nil {
		return nil, fmt.Errorf(
			"failed to execute get plant by id query for user_id=%s plant_id=%s: %w",
			userID, plantID, err,
		)
	}
	defer rows.Close()
	plantItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[plant.PopulatedPlant])
	if err != nil {
		return nil, fmt.Errorf(
			"failed to collect populated plant for user_id=%s plant_id=%s: %w",
			userID, plantID, err,
		)
	}

	return &plantItem, nil
}


func (p *PlantRepository) CheckPlantExists(ctx context.Context, userID string, plantID uuid.UUID) (*plant.Plant, error) {
	stmt := `
		SELECT
			*
		FROM
			plants
		WHERE
			id = @plant_id
			AND user_id = @user_id
	`

	rows, err := p.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"plant_id": plantID,
		"user_id":  userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check if plant exists for plant_id=%s user_id=%s: %w", plantID.String(), userID, err)
	}

	plantItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[plant.Plant])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:plants for plant_id=%s user_id=%s: %w", plantID.String(), userID, err)
	}

	return &plantItem, nil
}
func (p *PlantRepository) GetPlants(ctx context.Context,userID string,query *dto.GetPlantsQuery,) (*model.PaginatedResponse[plant.PopulatedPlant], error) {

	stmt := `
		SELECT
			p.*,

			-- Aggregate all observations for each plant
			COALESCE(
				jsonb_agg(
					to_jsonb(camel(obs))
					ORDER BY obs.date DESC
				) FILTER (WHERE obs.id IS NOT NULL),
				'[]'::JSONB
			) AS observations
		FROM plants p
		LEFT JOIN observations obs
			ON obs.plant_id = p.id
			AND obs.user_id = @user_id
	`

	args := pgx.NamedArgs{
		"user_id": userID,
	}
	conditions := []string{"p.user_id = @user_id"}

	// --- Dynamic filtering ---
	if query.Search != nil {
		conditions = append(conditions,
			"(p.name ILIKE @search OR p.species ILIKE @search OR p.location ILIKE @search)")
		args["search"] = "%" + *query.Search + "%"
	}

	if query.Species != nil {
		conditions = append(conditions, "p.species = @species")
		args["species"] = *query.Species
	}

	if query.Location != nil {
		conditions = append(conditions, "p.location = @location")
		args["location"] = *query.Location
	}

	if query.PlantedFrom != nil {
		conditions = append(conditions, "p.planted_date >= @planted_from")
		args["planted_from"] = *query.PlantedFrom
	}

	if query.PlantedTo != nil {
		conditions = append(conditions, "p.planted_date <= @planted_to")
		args["planted_to"] = *query.PlantedTo
	}

	if len(conditions) > 0 {
		stmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	// --- Count query for pagination ---
	countStmt := "SELECT COUNT(*) FROM plants p"
	if len(conditions) > 0 {
		countStmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := p.server.DB.Pool.QueryRow(ctx, countStmt, args).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count plants for user_id=%s: %w", userID, err)
	}

	// --- Grouping ---
	stmt += " GROUP BY p.id"

	// --- Sorting ---
	if query.Sort != nil {
		stmt += " ORDER BY p." + *query.Sort
		if query.Order != nil && *query.Order == "desc" {
			stmt += " DESC"
		} else {
			stmt += " ASC"
		}
	} else {
		stmt += " ORDER BY p.created_at DESC"
	}

	// --- Pagination ---
	stmt += " LIMIT @limit OFFSET @offset"
	args["limit"] = *query.Limit
	args["offset"] = (*query.Page - 1) * (*query.Limit)

	// --- Query execution ---
	rows, err := p.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get plants query for user_id=%s: %w", userID, err)
	}
	defer rows.Close()

	// --- Collect into structs ---
	plants, err := pgx.CollectRows(rows, pgx.RowToStructByName[plant.PopulatedPlant])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[plant.PopulatedPlant]{
				Data:       []plant.PopulatedPlant{},
				Page:       *query.Page,
				Limit:      *query.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to map plant rows for user_id=%s: %w", userID, err)
	}

	// --- Build pagination response ---
	return &model.PaginatedResponse[plant.PopulatedPlant]{
		Data:       plants,
		Page:       *query.Page,
		Limit:      *query.Limit,
		Total:      total,
		TotalPages: (total + *query.Limit - 1) / *query.Limit,
	}, nil
}

