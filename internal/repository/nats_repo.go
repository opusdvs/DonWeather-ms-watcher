package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/opusdvs/DonWeather-ms-watcher/internal/domain"
)

type EventRepository struct {
	conn *nats.Conn
}

func NewEventRepository(conn *nats.Conn) *EventRepository {
	return &EventRepository{conn: conn}
}

// нужно будет доработать, чтобы использовать контекст для отмены операции
func (r *EventRepository) CreateEventWeather(ctx context.Context, event *domain.EventWeather) error {
	message := domain.EventWeatherMessage{
		ID:        uuid.New().String(),
		Data:      *event,
		Timestamp: time.Now(),
		Type:      string(event.EventType),
		Source:    "weather-service",
		Version:   "1.0.0",
	}
	json, err := json.Marshal(message)
	if err != nil {
		return err
	}
	subject := "event.weather.message" + "." + string(event.EventType)
	err = r.conn.Publish(subject, json)
	if err != nil {
		return err
	}
	return nil
}
