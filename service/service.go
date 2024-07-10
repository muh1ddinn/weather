package service

import (
	"weather/pkg/logger"
	"weather/storage"
)

type IServiceMangaer interface {
	Weather() weatherserv
}

type Service struct {
	weather weatherserv

	logger logger.ILogger
}

func New(storage storage.IStorage, log logger.ILogger) Service {
	return Service{

		weather: NewCategories(storage, log),

		logger: log,
	}
}

func (s Service) Weather() weatherserv {
	return s.weather
}
