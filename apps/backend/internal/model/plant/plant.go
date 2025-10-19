package plant

import (
	"time"

	"github.com/google/uuid"
	"github.com/srauf24/gardenjournal/internal/model"
)
type Plant struct {
	model.Base                 // expects fields like ID, CreatedAt, UpdatedAt (and/or SortOrder if your Base includes it)
	UserID      string         `json:"userId" db:"user_id"`
	Name        string         `json:"name" db:"name"`
	Species     string         `json:"species" db:"species"`
	Location    *string        `json:"location" db:"location"`
	PlantedDate *time.Time     `json:"plantedDate" db:"planted_date"`
	Notes       *string        `json:"notes" db:"notes"`
	Metadata    json.RawMessage`json:"metadata" db:"metadata"`
	SortOrder   int            `json:"sortOrder" db:"sort_order"` // keep only if not already in model.Base
}
// Observation rows map 1:1 to the `observations` table.
type Observation struct {
	model.Base
	UserID   string     `json:"userId" db:"user_id"`
	PlantID  uuid.UUID  `json:"plantId" db:"plant_id"`
	Date     time.Time  `json:"date" db:"date"`
	HeightCM *float64   `json:"heightCm" db:"height_cm"`
	Notes    *string    `json:"notes" db:"notes"`
	SortOrder int       `json:"sortOrder" db:"sort_order"` // keep only if not in Base
}
// PopulatedPlant is the “rich” DTO for list/detail responses.
type PopulatedPlant struct {
	Plant
	Observations []Observation `json:"observations" db:"observations"`


// Metadata holds additional, optional properties about a plant.
// Stored as JSONB in the database, allowing flexible enrichment
// without requiring schema migrations.
type Metadata struct {
	Tags []string `json:"tags"` // e.g. ["indoor", "succulent"]

	// Care and Environment
	WateringFrequency *string    `json:"wateringFrequency"` // e.g. "every 3 days"
	LastWateredAt     *time.Time `json:"lastWateredAt"`
	SunlightLevel     *string    `json:"sunlightLevel"`  // e.g. "full sun", "partial shade"
	SoilType          *string    `json:"soilType"`       // e.g. "loamy", "sandy", "clay"
	PotSizeCM         *float64   `json:"potSizeCm"`      // e.g. 15.5
	FertilizerType    *string    `json:"fertilizerType"` // e.g. "NPK 10-10-10"
	LastFertilizedAt  *time.Time `json:"lastFertilizedAt"`

	// Environmental Tracking (weather integration)
	LastWeatherSnapshotID *uuid.UUID `json:"lastWeatherSnapshotId"`
	AverageTempC          *float64   `json:"averageTempC"`
	AverageSunshineHrs    *float64   `json:"averageSunshineHrs"`

	// Growth and Health
	HealthStatus *string  `json:"healthStatus"` // e.g. "thriving", "wilting"
	GrowthStage  *string  `json:"growthStage"`  // e.g. "seedling", "blooming"
	HeightCM     *float64 `json:"heightCm"`     // most recent height reading

	// Visual / UI Enrichment
	ColorTag  *string `json:"colorTag"`  // e.g. "#4CAF50" for display in UI
	ImageURL  *string `json:"imageUrl"`  // optional link to uploaded plant image
	EmojiIcon *string `json:"emojiIcon"` // e.g. "🌵", "🌿", "🌸"

	// Future Expansion (AI Insights, Summaries)
	AIInsightSummary *string `json:"aiInsightSummary"`
}


//Future: add a populated Plant DTO when the UI starts needing richer, nested data (e.g., observations with weather snapshots).