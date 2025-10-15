package observation


import (
	"time"
	"github.com/google/uuid"
	"github.com/srauf24/gardenjournal/internal/model"
)

// Observation represents a growth record for a plant.
type WeatherSnapshot struct {
    ID          uuid.UUID `json:"id" db:"id"`
    Date        *time.Time `json:"date" db:"date"`
    City        string    `json:"city" db:"city"`
    Latitude    Latitude   `json:"latitude" db:"latitude"`
    Longitude   Longitude   `json:"longitude" db:"longitude"`
    TempMax     *float64  `json:"tempMax" db:"temp_max"`
    PrecipMM    *float64  `json:"precipMm" db:"precip_mm"`
    SunshineHrs *float64  `json:"sunshineHrs" db:"sunshine_hrs"`
}


