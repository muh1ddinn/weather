package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"weather/api/model"

	"github.com/gin-gonic/gin"
)

// GetById godoc
// @Security ApiKeyAuth
// @Router     /weather [GET]
// @Summary    Get weather data by country name
// @Description This API gets weather data by country name and returns its info
// @Tags       weather
// @Accept     json
// @Produce    json
// @Param      search query string false "city_name"
// @Success    200  {object}  model.WeatherResponse
// @Failure    400  {object}  model.Response
// @Failure    404  {object}  model.Response
// @Failure    500  {object}  model.Response
func (h Handler) Getweather(c *gin.Context) {
	country := c.Query("search")
	fmt.Println("search ", country)

	// Get location by country name
	location, err := h.Services.Weather().Country(context.Background(), country)
	if err != nil {
		handleResponseLog(c, h.Log, "error while getting location by country name", http.StatusBadRequest, err.Error())
		return
	}

	// Check if location ID is empty (indicating no rows found)
	if location.Id != "" {
		// Delete existing weather data for the location
		err = h.Services.Weather().Delete(context.Background(), location)
		if err != nil {
			handleResponseLog(c, h.Log, "error while deleting location and weather data", http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Retrieve new weather data by country name
	weather, err := h.Services.Weather().Getahour(context.Background(), country)
	if err != nil {
		handleResponseLog(c, h.Log, "error while getting weather data by country name", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "weather data retrieved successfully by country name", http.StatusOK, weather)
}

// GetAllWeather godoc
// @Security ApiKeyAuth
// @Router     /weatherget [GET]
// @Summary    Get all weather data with filters
// @Description This API gets all weather data with optional filters
// @Tags       weather
// @Accept     json
// @Produce    json
// @Param      city query string false "City name"
// @Param      condition query string false "Weather condition"
// @Param      min_temp query float64 false "Minimum temperature"
// @Param      max_temp query float64 false "Maximum temperature"
// @Param      start_time query string false "Start time (YYYY-MM-DD)"
// @Param      end_time query string false "End time (YYYY-MM-DD)"
// @Param      page query int false "Page number"
// @Param      limit query int false "Limit"
// @Success    200 {object} model.GetAllWeatherResponse
// @Failure    400 {object} model.Response
// @Failure    404 {object} model.Response
// @Failure    500 {object} model.Response
func (h Handler) GetAllWeather(c *gin.Context) {
	var req model.GetAllWeatherRequestt

	// Parse query parameters into the req struct
	req.City = c.Query("city")
	req.Condition = c.Query("condition")
	req.StartTime = c.Query("start_time")
	req.EndTime = c.Query("end_time")

	// Parse min_temp and max_temp parameters
	if minTempStr := c.Query("min_temp"); minTempStr != "" {
		minTemp, err := strconv.ParseFloat(minTempStr, 64)
		if err != nil {
			handleResponseLog(c, h.Log, "error while parsing min_temp", http.StatusBadRequest, err.Error())
			return
		}
		req.MinTemp = minTemp
	}

	if maxTempStr := c.Query("max_temp"); maxTempStr != "" {
		maxTemp, err := strconv.ParseFloat(maxTempStr, 64)
		if err != nil {
			handleResponseLog(c, h.Log, "error while parsing max_temp", http.StatusBadRequest, err.Error())
			return
		}
		req.MaxTemp = maxTemp
	}

	// Parse page and limit parameters
	page, err := ParsePageQueryParam(c)
	if err != nil {
		handleResponseLog(c, h.Log, "error while parsing page", http.StatusBadRequest, err.Error())
		return
	}
	req.Page = (page)

	limit, err := ParseLimitQueryParam(c)
	if err != nil {
		handleResponseLog(c, h.Log, "error while parsing limit", http.StatusBadRequest, err.Error())
		return
	}
	req.Limit = (limit)

	// Ensure default values for page and limit
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	fmt.Println("city:", req.City)
	fmt.Println("city:", req.MaxTemp)

	fmt.Println("city:", req.StartTime)

	// Call service method to fetch weather data
	weather, err := h.Services.Weather().GetAllWeather(context.Background(), req)
	if err != nil {
		handleResponseLog(c, h.Log, "error while getting weather data", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "weather data retrieved successfully", http.StatusOK, weather)
}
