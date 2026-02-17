package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/opusdvs/DonWeather-ms-watcher/internal/domain"
	"github.com/redis/go-redis/v9"
)

type WeatherRequest struct {
	City string `json:"q"`
	Days string `json:"days"`
	Lang string `json:"lang"`
}
type WeatherRepository struct {
	redisClient   *redis.Client
	httpClient    *http.Client
	weatherApiURL string
}

func NewWeatherRepository(redisClient *redis.Client, weatherApiURL string) *WeatherRepository {
	return &WeatherRepository{
		redisClient:   redisClient,
		httpClient:    &http.Client{Timeout: 10 * time.Second},
		weatherApiURL: weatherApiURL,
	}
}

func (r *WeatherRepository) SaveStateWeather(ctx context.Context, city string, stateWeather *domain.StateWeather) error {

	data, err := json.Marshal(stateWeather)
	if err != nil {
		return err
	}
	return r.redisClient.Set(ctx, city, data, 0).Err()
}

func (r *WeatherRepository) GetLastStateWeather(ctx context.Context, city string) (*domain.StateWeather, error) {
	log.Println("Getting last state weather for city:", city)
	data, err := r.redisClient.Get(ctx, city).Bytes()
	if err != nil {
		return nil, err
	}
	var stateWeather domain.StateWeather
	if err = json.Unmarshal(data, &stateWeather); err != nil {
		return nil, err
	}
	return &stateWeather, nil
}

func (r *WeatherRepository) GetCurrentStateWeather(ctx context.Context, city string) (*domain.StateWeather, error) {
	var stateWeather domain.StateWeather
	request := WeatherRequest{
		City: city,
		Days: "3",
		Lang: "ru",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.weatherApiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get current state weather: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &stateWeather)
	if err != nil {
		return nil, err
	}
	return &stateWeather, nil
}
