package plant

import (
	"time"

	"github.com/google/uuid"
	"github.com/srauf24/gardenjournal/internal/model"
)

// Plant represents a single tracked plant record.
// Nullable values are represented with pointer types.
//   - `json` tags define how the struct is serialized when communicating with the frontend/API
//   - `db` tags define how the struct maps to the database columns in PostgreSQL
//     (used by the SQL driver for reading/writing data)
type Plant struct {
	model.Base
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      string     `json:"userID" db:"user_id"`
	Name        Name       `json:"name" db:"name"`
	Species     Species    `json:"species" db:"species"`
	Location    *string    `json:"location" db:"location"`
	PlantedDate *time.Time `json:"plantedDate" db:"planted_date"`
	Notes       *string    `json:"notes" db:"notes"`
	Metadata    *Metadata  `json:"metadata" db:"metadata"`
	SortOrder   int        `json:"sortOrder" db:"sort_order"`
}

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
