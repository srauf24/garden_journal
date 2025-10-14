CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- --- Plants Table ---
CREATE TABLE plants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    user_id TEXT NOT NULL,                 -- Foreign key to the user who owns the plant.
    name TEXT NOT NULL,                    -- The name given to the plant by the user.
    species TEXT NOT NULL,                 -- The species of the plant.
    location TEXT,                         -- Where the plant is located (e.g., "living room", "garden").
    planted_date TIMESTAMPTZ,              -- The date the plant was planted.
    notes TEXT                             -- Optional field for general notes.
);

-- Indexes for plants
CREATE INDEX idx_plants_user_id ON plants(user_id);
CREATE UNIQUE INDEX plants_unique_name_per_user ON plants(user_id, name);

-- Trigger for plants
CREATE TRIGGER set_updated_at_plants
    BEFORE UPDATE ON plants
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- --- Weather Snapshots Table ---
CREATE TABLE weather_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    date TIMESTAMPTZ NOT NULL,             -- The date of the weather snapshot.
    city TEXT NOT NULL,                    -- The city for which the weather data was recorded.
    latitude DOUBLE PRECISION NOT NULL,    -- The latitude of the location.
    longitude DOUBLE PRECISION NOT NULL,   -- The longitude of the location.
    temp_max REAL,                         -- The maximum temperature for the day.
    precip_mm REAL,                        -- The amount of precipitation in millimeters.
    sunshine_hrs REAL                      -- The number of hours of sunshine.
);

-- Indexes for weather_snapshots
CREATE INDEX idx_weather_snapshots_city_date ON weather_snapshots(city, date);
CREATE INDEX idx_weather_snapshots_lat_long ON weather_snapshots(latitude, longitude);

-- --- Observations Table ---
CREATE TABLE observations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    user_id TEXT NOT NULL,                                   -- Foreign key to the user who made the observation.
    plant_id UUID NOT NULL REFERENCES plants(id) ON DELETE CASCADE,  -- The plant being observed.
    weather_snapshot_id UUID REFERENCES weather_snapshots(id) ON DELETE SET NULL,  -- Optional link to weather.
    date TIMESTAMPTZ NOT NULL,                               -- The date of the observation.
    height_cm REAL,                                          -- The height of the plant in centimeters.
    notes TEXT                                               -- Any notes about the observation.
);

-- Indexes for observations
CREATE INDEX idx_observations_user_id ON observations(user_id);
CREATE INDEX idx_observations_plant_date ON observations(plant_id, date);
CREATE INDEX idx_observations_weather_snapshot_id ON observations(weather_snapshot_id);

-- Trigger for observations
CREATE TRIGGER set_updated_at_observations
    BEFORE UPDATE ON observations
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- --- Constraints & Enhancements ---
-- Prevent duplicate observations for the same plant on the same date
ALTER TABLE observations
ADD CONSTRAINT unique_plant_observation_per_day UNIQUE (plant_id, date);

-- Optional composite index for quick retrieval by user/date
CREATE INDEX idx_observations_user_date ON observations(user_id, date);
