package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/srauf24/gardenjournal//internal/model/observation"
	"github.com/srauf24/gardenjournal/internal/server"
)

type ObservationRepository struct {
    server *server.Server
}
func NewObservationRepository(server *server.Server) *ObservationRepository {
    return &ObservationRepository{server: server}
}
// note: typically the service layer translates the DTO to a model which is then used by the repository layer. In this case we are directly using the DTO in the repository layer for simplicity.
func (r *ObservationRepository) CreateObservation(ctx context.Context, userID string, payload *observation.CreateObservationPayload) (*observation.Observation, error) {
    stmt := `
        INSERT INTO
            observations (
                user_id,
                plant_id,
                date,
                height_cm,
                notes,
            )
        VALUES
            (
                @user_id,
                @plant_id,
                @date,
                @height_cm,
                @notes,
            )
        RETURNING
        *
    `
    // use server.db to execute the query
    rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
        "user_id":        userID,
        "plant_id":      payload.PlantID,
        "date":          payload.Date,
        "height_cm":    payload.HeightCM,
        "notes":         payload.Notes,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to execute create observation query for user_id=%s plant_id=%s: %w", userID, payload.PlantID, err)
    }
    // use pgx library to deserialize row into a struct from the data base (collect one row)
    observationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[observation.Observation])
    if err != nil {
        return nil, fmt.Errorf("failed to collect row from table:observations for user_id=%s plant_id=%s: %w", userID, payload.PlantID, err)
    }
    return observationItem, nil
}