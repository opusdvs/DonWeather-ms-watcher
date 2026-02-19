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

func (r *WeatherRepository) SaveStateWeather(ctx context.Context, city string, stateWeather domain.StateWeather) error {
	// Устанавливаем поля перед сохранением
	stateWeather.Location.Name = city
	if stateWeather.UpdatedAt.IsZero() {
		stateWeather.UpdatedAt = time.Now()
	}

	// Сериализуем структуру в JSON
	data, err := json.Marshal(stateWeather)
	if err != nil {
		log.Println("Failed to marshal stateWeather:", err)
		return err
	}

	// Сохраняем в Redis (ключ = city, без TTL)
	err = r.redisClient.Set(ctx, city, data, 0).Err()
	if err != nil {
		log.Println("Redis SET error:", err)
		return err
	}

	log.Printf("Weather state saved successfully for city: %s (temp=%.2f, humidity=%.2f, wind=%.2f, pressure=%.2f)",
		city, stateWeather.Current.TempC, stateWeather.Current.Humidity, stateWeather.Current.WindKph, stateWeather.Current.PressureMb)
	return nil
}

func (r *WeatherRepository) GetLastStateWeather(ctx context.Context, city string) (domain.StateWeather, error) {
	log.Println("Getting last state weather for city:", city)
	data, err := r.redisClient.Get(ctx, city).Bytes()
	if err != nil {
		return domain.StateWeather{}, err
	}
	var stateWeather domain.StateWeather
	if err = json.Unmarshal(data, &stateWeather); err != nil {
		return domain.StateWeather{}, err
	}
	return stateWeather, nil
}

func (r *WeatherRepository) GetCurrentStateWeather(ctx context.Context, city string) (domain.StateWeather, error) {
	var stateWeather domain.StateWeather
	request := WeatherRequest{
		City: city,
		Days: "3",
		Lang: "ru",
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return domain.StateWeather{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("Requesting weather for city: %s, URL: %s", city, r.weatherApiURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.weatherApiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return domain.StateWeather{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return domain.StateWeather{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.StateWeather{}, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Weather API response status: %d, body length: %d", resp.StatusCode, len(body))
	if resp.StatusCode != http.StatusOK {
		log.Printf("Weather API error response: %s", string(body))
		return domain.StateWeather{}, fmt.Errorf("failed to get current state weather: status %d, %s", resp.StatusCode, resp.Status)
	}

	log.Printf("Weather API response body: %s", string(body))

	err = json.Unmarshal(body, &stateWeather)
	if err != nil {
		log.Printf("Failed to unmarshal weather response: %v, body: %s", err, string(body))
		return domain.StateWeather{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Устанавливаем поля, которые могут отсутствовать в ответе API
	stateWeather.Location.Name = city
	stateWeather.UpdatedAt = time.Now()

	log.Printf("Parsed weather data for %s: temp=%.2f, humidity=%.2f, wind=%.2f, pressure=%.2f",
		city, stateWeather.Current.TempC, stateWeather.Current.Humidity, stateWeather.Current.WindKph, stateWeather.Current.PressureMb)

	return stateWeather, nil
}
