package repository

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/opusdvs/DonWeather-ms-watcher/internal/domain"
)

type EventRepository struct {
	natsConn nats.JetStreamContext
}

func NewEventRepository(natsConn *nats.Conn) (*EventRepository, error) {
	js, err := natsConn.JetStream()
	if err != nil {
		return nil, err
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "WEATHER_EVENTS",
		Subjects:  []string{"event.weather.message.*"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
	})
	if err != nil {
		return nil, err
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "TELEGRAM_EVENTS",
		Subjects:  []string{"event.telegram.message.send"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
	})
	if err != nil {
		return nil, err
	}
	return &EventRepository{natsConn: js}, nil
}

// нужно будет доработать, чтобы использовать контекст для отмены операции
func (r *EventRepository) CreateEventWeather(ctx context.Context, subscription domain.Subscribe, event domain.EventWeather) error {
	message := domain.EventWeatherMessage{
		ID:         uuid.New().String(),
		TelegramID: subscription.TelegramID,
		Data:       event,
		Timestamp:  time.Now(),
		Type:       string(event.EventType),
		Source:     "weather-service",
		Version:    "1.0.0",
	}
	log.Println("Event:", message)
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	subject := "event.weather.message" + "." + string(event.EventType)
	_, err = r.natsConn.Publish(subject, jsonData)
	if err != nil {
		return err
	}
	return nil
}
