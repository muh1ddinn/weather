package storage

import (
	"context"
	"weather/api/model"
)

type IStorage interface {
	CloseDB()
	Weather() IWeatherStorage
}

type IWeatherStorage interface {
	GetAllWeather(context.Context, model.GetAllWeatherRequestt) (model.GetAllWeatherResponse, error)
	Get(context.Context, string) (model.WeatherResponse, error)
	Delete(context.Context, model.GetAllWeathername) error
	Country(context.Context, string) (model.GetAllWeathername, error)
	UpdateWeatherData(ctx context.Context, location string) error
}
