package repository

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/opusdvs/DonWeather-ms-watcher/internal/domain"
)

type EventRepository struct {
	conn *nats.Conn
}

func NewEventRepository(conn *nats.Conn) *EventRepository {
	return &EventRepository{conn: conn}
}

func (r *EventRepository) CreateEventWeather(ctx context.Context, event *domain.EventWeather) error {
	json, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = r.conn.Publish(string(event.EventType), json)
	if err != nil {
		return err
	}
	return nil
}
