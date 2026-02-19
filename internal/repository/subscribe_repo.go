package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/opusdvs/DonWeather-ms-watcher/internal/domain"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) GetActiveSubscriptions(ctx context.Context) ([]domain.Subscribe, error) {
	query := "SELECT id, city, filters, telegram_id FROM subscribe"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	subscriptions := []domain.Subscribe{}
	for rows.Next() {
		var subscription domain.Subscribe
		var rawFilters []byte
		var telegramID string
		err = rows.Scan(&subscription.ID, &subscription.City, &rawFilters, &telegramID)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(rawFilters, &subscription.Filters)
		if err != nil {
			return nil, err
		}
		subscription.TelegramID = telegramID
		subscriptions = append(subscriptions, subscription)
		log.Println("Subscription:", subscriptions)
	}
	return subscriptions, nil
}
