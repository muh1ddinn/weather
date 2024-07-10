
-- Create the locations table
CREATE TABLE locations (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL
);

-- Create the weather_hour table
CREATE TABLE weather_hour (
    id SERIAL PRIMARY KEY,
    time TIMESTAMP NOT NULL,
    temperature FLOAT NOT NULL,
    humidity FLOAT NOT NULL,
    condition VARCHAR(100) NOT NULL,
    location_id UUID REFERENCES locations(id)
);

-- Create the weather_day table
CREATE TABLE weather_day (
    id SERIAL PRIMARY KEY,
    time TIMESTAMP NOT NULL,
    temperature FLOAT NOT NULL,
    temperaturemax FLOAT NOT NULL,
    temperaturemin FLOAT NOT NULL,
    humidity FLOAT NOT NULL,
    condition VARCHAR(100) NOT NULL,
    location_id UUID REFERENCES locations(id)
);
