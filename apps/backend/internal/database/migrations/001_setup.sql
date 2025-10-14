CREATE TABLE plants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    user_id TEXT NOT NULL, -- Foreign key to the user who owns the plant.
    name TEXT NOT NULL, -- The name given to the plant by the user.
    species TEXT NOT NULL, -- The species of the plant.
    location TEXT, -- Where the plant is located (e.g., "living room", "garden").
    planted_date TIMESTAMPTZ -- The date the plant was planted.
);
CREATE INDEX idx_plants_user_id ON plants(user_id);

CREATE TABLE weather_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    date TIMESTAMPTZ NOT NULL, -- The date of the weather snapshot.
    city TEXT NOT NULL, -- The city for which the weather data was recorded.
    latitude DOUBLE PRECISION NOT NULL, -- The latitude of the location.
    longitude DOUBLE PRECISION NOT NULL, -- The longitude of the location.
    temp_max REAL, -- The maximum temperature for the day.
    precip_mm REAL, -- The amount of precipitation in millimeters.
    sunshine_hrs REAL -- The number of hours of sunshine.
);
CREATE INDEX idx_weather_snapshots_city_date ON weather_snapshots(city, date);

CREATE TABLE observations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    user_id TEXT NOT NULL, -- Foreign key to the user who made the observation.
    plant_id UUID NOT NULL REFERENCES plants(id) ON DELETE CASCADE, -- The plant being observed.
    weather_snapshot_id UUID REFERENCES weather_snapshots(id) ON DELETE SET NULL, -- Optional link to a weather snapshot.
    date TIMESTAMPTZ NOT NULL, -- The date of the observation.
    height_cm REAL, -- The height of the plant in centimeters.
    notes TEXT -- Any notes about the observation.
);
CREATE INDEX idx_observations_plant_date ON observations(plant_id, date);


RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Attach triggers to relevant tables.
CREATE TRIGGER set_updated_at_plants
BEFORE UPDATE ON plants
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at_observations
BEFORE UPDATE ON observations
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();
