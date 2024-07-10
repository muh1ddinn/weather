package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"weather/api/model"
	"weather/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type WeatherRepo struct {
	db     *pgxpool.Pool
	logger logger.ILogger
}

// NewWeatherRepo creates a new instance of WeatherRepo
func NewWeatherRepo(db *pgxpool.Pool, log logger.ILogger) WeatherRepo {
	return WeatherRepo{
		db:     db,
		logger: log,
	}
}

func (c *WeatherRepo) Get(ctx context.Context, location string) (model.WeatherResponse, error) {
	hourlyResponse := model.WeatherResponse{}
	apiKey := "FK56GRHNE2E9UQQPK8KVP3NUW"
	baseURL := "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline"

	url := fmt.Sprintf("%s/%s?key=%s", baseURL, location, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return hourlyResponse, fmt.Errorf("error making the request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return hourlyResponse, fmt.Errorf("error reading the response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return hourlyResponse, fmt.Errorf("received non-200 response status: %s, response body: %s", resp.Status, string(body))
	}

	var weatherResponse model.WeatherResponse
	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		return hourlyResponse, fmt.Errorf("error parsing the JSON response: %w", err)
	}

	// Print the general weather data
	fmt.Printf("Location: %s\n", weatherResponse.Address)
	fmt.Printf("Description: %s\n", weatherResponse.ResolvedAddress)

	// Insert location into the database
	locationID := uuid.New()
	_, err = c.db.Exec(ctx, "INSERT INTO locations (id, name, country) VALUES ($1, $2, $3)",
		locationID, weatherResponse.Address, weatherResponse.ResolvedAddress)
	if err != nil {
		return hourlyResponse, fmt.Errorf("error inserting location into DB: %w", err)
	}

	for _, hour := range weatherResponse.Days[0].Hours {
		hourTime := time.Unix(hour.DatetimeEpoch, 0)
		_, err = c.db.Exec(ctx, "INSERT INTO weather_hour (time, temperature, humidity, condition, location_id) VALUES ($1, $2, $3, $4, $5)",
			hourTime, hour.Temp, hour.Humidity, hour.Icon, locationID)
		if err != nil {
			return hourlyResponse, fmt.Errorf("error inserting hourly data into DB: %w", err)
		}
	}

	for _, day := range weatherResponse.Days {
		dayTime := time.Unix(day.DatetimeEpoch, 0)
		_, err = c.db.Exec(ctx, "INSERT INTO weather_day (time, temperature, temperaturemax, temperaturemin, humidity, condition, location_id) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			dayTime, day.Temp, day.TempMax, day.TempMin, day.Humidity, day.Conditions, locationID)
		if err != nil {
			return hourlyResponse, fmt.Errorf("error inserting daily data into DB: %w", err)
		}
	}

	hourlyResponse.Hours = append(hourlyResponse.Hours, weatherResponse.Hours...)

	hourlyResponse.Days = append(hourlyResponse.Days, weatherResponse.Days...)

	for _, day := range weatherResponse.Days {
		date := time.Unix(day.DatetimeEpoch, 0)
		fmt.Printf(
			"%s - %.1f°F,%.1f°F,%.1f°F, %.1f%% humidity, %s\n",
			date.Format("2006-01-02 15:04"),
			day.Temp,
			day.TempMax,
			day.TempMin,
			day.Humidity,
			day.Conditions,
		)

	}
	return hourlyResponse, nil
}
func (c *WeatherRepo) Delete(ctx context.Context, del model.GetAllWeathername) error {
	fmt.Println(del.Id, "delprint")

	query := `DELETE FROM weather_hour WHERE location_id = $1`
	_, err := c.db.Exec(ctx, query, del.Id)
	if err != nil {
		c.logger.Error("failed to delete weather_hour from database", logger.Error(err))
		fmt.Println(err, "errdelet")
		return err
	}

	query = `DELETE FROM weather_day WHERE location_id = $1`
	_, err = c.db.Exec(ctx, query, del.Id)
	if err != nil {
		c.logger.Error("failed to delete weather_day from database", logger.Error(err))
		return err
	}

	query = `DELETE FROM locations WHERE id = $1`
	_, err = c.db.Exec(ctx, query, del.Id)
	if err != nil {
		c.logger.Error("failed to delete location from database", logger.Error(err))
		return err
	}

	return nil
}
func (c *WeatherRepo) Country(ctx context.Context, name string) (model.GetAllWeathername, error) {
	locations := model.GetAllWeathername{}
	query := `
        SELECT id, name
        FROM locations 
        WHERE name=$1 
    `

	row := c.db.QueryRow(ctx, query, name)

	err := row.Scan(&locations.Id, &locations.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			return locations, nil // No error, but return empty locations
		}
		fmt.Println(err, "errpostgrecountryname")
		return locations, err
	}

	return locations, nil
}

func (c *WeatherRepo) GetAllWeather(ctx context.Context, req model.GetAllWeatherRequestt) (model.GetAllWeatherResponse, error) {
	fmt.Println("citysearch", req.City)
	fmt.Println("citysearch", req.MaxTemp)

	resp := model.GetAllWeatherResponse{}
	offset := (req.Page - 1) * req.Limit
	filter := ""
	args := []interface{}{}
	argIndex := 1

	// Adding city name filter
	if req.City != "" {
		filter += fmt.Sprintf(` AND l.name ILIKE $%d `, argIndex)
		args = append(args, "%"+req.City+"%")
		argIndex++
	}

	// Adding condition filter
	if req.Condition != "" {
		filter += fmt.Sprintf(` AND (wh.condition ILIKE $%d OR wd.condition ILIKE $%d) `, argIndex, argIndex+1)
		args = append(args, "%"+req.Condition+"%", "%"+req.Condition+"%")
		argIndex += 2
	}

	// Adding temperature range filter
	if req.MinTemp != 0 && req.MaxTemp != 0 {
		filter += fmt.Sprintf(` AND ((wh.temperature BETWEEN $%d AND $%d) OR (wd.temperaturemin BETWEEN $%d AND $%d) OR (wd.temperaturemax BETWEEN $%d AND $%d)) `, argIndex, argIndex+1, argIndex+2, argIndex+3, argIndex+4, argIndex+5)
		args = append(args, req.MinTemp, req.MaxTemp, req.MinTemp, req.MaxTemp, req.MinTemp, req.MaxTemp)
		argIndex += 6
	} else if req.MinTemp != 0 {
		filter += fmt.Sprintf(` AND (wh.temperature >= $%d OR wd.temperaturemin >= $%d) `, argIndex, argIndex+1)
		args = append(args, req.MinTemp, req.MinTemp)
		argIndex += 2
	} else if req.MaxTemp != 0 {
		filter += fmt.Sprintf(` AND (wh.temperature <= $%d OR wd.temperaturemax <= $%d) `, argIndex, argIndex+1)
		args = append(args, req.MaxTemp, req.MaxTemp)
		argIndex += 2
	}

	// Adding start time filter
	if req.StartTime != "" {
		startTime, err := time.Parse("2006-01-02", req.StartTime)
		if err != nil {
			return resp, fmt.Errorf("invalid start_time format")
		}
		filter += fmt.Sprintf(` AND (wh.time >= $%d OR wd.time <= $%d) `, argIndex, argIndex+1)
		args = append(args, startTime, startTime)
		argIndex += 2
	}

	// Adding end time filter
	if req.EndTime != "" {
		endTime, err := time.Parse("2006-01-02", req.EndTime)
		if err != nil {
			return resp, fmt.Errorf("invalid end_time format")
		}
		filter += fmt.Sprintf(` AND (wh.time <= $%d OR wd.time <= $%d) `, argIndex, argIndex+1)
		args = append(args, endTime, endTime)
		argIndex += 2
	}

	// Construct the main query
	query := `
        SELECT 
            l.id,
            l.name,
            l.country,
            wh.time AS hour_time,
            wh.temperature AS hour_temperature,
            wh.humidity AS hour_humidity,
            wh.condition AS hour_condition,
            wd.time AS day_time,
            wd.temperature AS day_temperature,
            wd.temperaturemax AS day_temperaturemax,
            wd.temperaturemin AS day_temperaturemin,
            wd.humidity AS day_humidity,
            wd.condition AS day_condition,
            count(*) OVER() AS total_count
        FROM 
            locations l
        LEFT JOIN
            weather_hour wh ON l.id = wh.location_id
        LEFT JOIN
            weather_day wd ON l.id = wd.location_id
        WHERE 1=1 ` + filter +
		fmt.Sprintf(" ORDER BY l.name ASC OFFSET $%d LIMIT $%d", argIndex, argIndex+1)

	args = append(args, offset, req.Limit)

	// Debugging output
	// fmt.Println("Constructed SQL query:", query)
	// fmt.Println("Arguments:", args)

	// Execute the query
	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	var totalCount int
	for rows.Next() {
		var weather model.Weather
		var hourTime, dayTime pq.NullTime // Using pq.NullTime for nullable time fields

		err := rows.Scan(
			&weather.Location.Id,
			&weather.Location.Name,
			&weather.Location.Country,
			&hourTime,
			&weather.Hour.Temperature,
			&weather.Hour.Humidity,
			&weather.Hour.Condition,
			&dayTime,
			&weather.Day.Temperature,
			&weather.Day.TemperatureMax,
			&weather.Day.TemperatureMin,
			&weather.Day.Humidity,
			&weather.Day.Condition,
			&totalCount,
		)
		if err != nil {
			fmt.Println(err, "err")
			return resp, err
		}

		// Convert time to string if valid
		if hourTime.Valid {
			weather.Hour.Time = hourTime.Time.Format("2006-01-02 15:04:05")
		}
		if dayTime.Valid {
			weather.Day.Time = dayTime.Time.Format("2006-01-02")
		}

		resp.Weathers = append(resp.Weathers, weather)
	}

	resp.TotalCount = totalCount
	return resp, nil
}

func (c *WeatherRepo) UpdateWeatherData(ctx context.Context, location string) error {
	weatherResponse, err := c.Get(ctx, location)
	if err != nil {
		return fmt.Errorf("failed to fetch weather data: %w", err)
	}

	locationID, err := c.getLocationID(ctx, weatherResponse.Address)
	if err != nil {
		return fmt.Errorf("failed to get location ID: %w", err)
	}

	for _, hour := range weatherResponse.Days[0].Hours {
		hourTime := time.Unix(hour.DatetimeEpoch, 0)
		_, err = c.db.Exec(ctx, "UPDATE weather_hour SET temperature = $1, humidity = $2, condition = $3 WHERE time = $4 AND location_id = $5",
			hour.Temp, hour.Humidity, hour.Icon, hourTime, locationID)
		if err != nil {
			return fmt.Errorf("error updating hourly data in DB: %w", err)
		}
	}

	for _, day := range weatherResponse.Days {
		dayTime := time.Unix(day.DatetimeEpoch, 0)
		_, err = c.db.Exec(ctx, "UPDATE weather_day SET temperature = $1, temperaturemax = $2, temperaturemin = $3, humidity = $4, condition = $5 WHERE time = $6 AND location_id = $7",
			day.Temp, day.TempMax, day.TempMin, day.Humidity, day.Conditions, dayTime, locationID)
		if err != nil {
			return fmt.Errorf("error updating daily data in DB: %w", err)
		}
	}

	return nil
}

func (c *WeatherRepo) getLocationID(ctx context.Context, location string) (uuid.UUID, error) {
	var locationID uuid.UUID
	err := c.db.QueryRow(ctx, "SELECT id FROM locations WHERE name = $1", location).Scan(&locationID)
	if err != nil {
		return uuid.Nil, err
	}
	return locationID, nil
}
