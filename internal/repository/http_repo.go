package repository

import (
	"context"

	"github.com/opusdvs/DonWeather-ms-watcher/internal/domain"
	"github.com/redis/go-redis/v9"
)

type WeatherRepository struct {
	redisClient *redis.Client
}

func NewWeatherRepository(redisClient *redis.Client) *WeatherRepository {
	return &WeatherRepository{redisClient: redisClient}
}

func (r *WeatherRepository) GetLastStateWeather(ctx context.Context, city string) (*domain.StateWeather, error) {
	return nil, nil
}

func (r *WeatherRepository) UpdateStateWeather(ctx context.Context, state *domain.StateWeather) error {
	return nil
}
