package service

import (
	"context"
	"fmt"
	"weather/api/model"
	"weather/pkg/logger"
	"weather/storage"
)

type weatherserv struct {
	storage storage.IStorage
	logger  logger.ILogger
}

func NewCategories(storage storage.IStorage, logger logger.ILogger) weatherserv {

	return weatherserv{
		storage: storage,
		logger:  logger,
	}
}

func (u weatherserv) Getahour(ctx context.Context, city string) (model.WeatherResponse, error) {

	msg, err := u.storage.Weather().Get(ctx, city)
	if err != nil {
		fmt.Println(err, "errservice")
		u.logger.Error("ERROR in service layer while create :waether", logger.Error(err))
		return model.WeatherResponse{}, err
	}

	return msg, nil
}
func (u weatherserv) Country(ctx context.Context, location string) (model.GetAllWeathername, error) {

	msg, err := u.storage.Weather().Country(ctx, location)
	if err != nil {
		fmt.Println(err, "errservice")
		u.logger.Error("ERROR in service layer while create :waether", logger.Error(err))
		return model.GetAllWeathername{}, err
	}

	return msg, nil
}

func (u weatherserv) Delete(ctx context.Context, del model.GetAllWeathername) error {

	err := u.storage.Weather().Delete(ctx, del)
	if err != nil {
		fmt.Println(err, "errservice")
		u.logger.Error("ERROR in service layer while create :waether", logger.Error(err))
		return err
	}

	return nil
}

func (u weatherserv) GetAllWeather(ctx context.Context, req model.GetAllWeatherRequestt) (model.GetAllWeatherResponse, error) {

	fmt.Println("req.time service", req.StartTime)

	msg, err := u.storage.Weather().GetAllWeather(ctx, req)
	fmt.Println(msg, "service")
	if err != nil {
		fmt.Println(err, "errservice")
		u.logger.Error("ERROR in service layer while create :waether", logger.Error(err))
		return model.GetAllWeatherResponse{}, err
	}

	return msg, nil

}

func (u weatherserv) Weatherupdte4(ctx context.Context, location string) error {

	err := u.storage.Weather().UpdateWeatherData(ctx, location)
	if err != nil {
		fmt.Println(err, "errservice")
		u.logger.Error("ERROR in service layer while create :update", logger.Error(err))
		return err
	}

	return nil

}
