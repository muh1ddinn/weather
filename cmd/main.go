package main

import (
	"context"
	"fmt"
	"time"
	"weather/api"
	"weather/config"
	"weather/pkg/logger"
	"weather/service"
	"weather/storage/postgres"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)

	store, err := postgres.New(context.Background(), cfg, log)
	if err != nil {
		fmt.Println("error while connecting db, err: ", err)
		return
	}
	defer store.CloseDB()

	service := service.New(store, log)

	// Start the periodic update
	go startPeriodicUpdate(context.Background(), &service, "Tashkent")

	c := api.New(service, log)

	fmt.Println("program is running on localhost:9090...")
	c.Run(":9091")
}

func startPeriodicUpdate(ctx context.Context, service *service.Service, location string) {
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := service.Weather().Weatherupdte4(ctx, location)
			if err != nil {
				fmt.Printf("Failed to update weather data: %v\n", err)
			} else {
				fmt.Println("Weather data updated successfully")
			}
		case <-ctx.Done():
			fmt.Println("Shutting down periodic update")
			return
		}
	}
}
