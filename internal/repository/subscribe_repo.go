package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/opusdvs/DonWeather-ms-watcher/internal/domain"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) GetActiveSubscriptions(ctx context.Context) ([]domain.Subscribe, error) {
	query := "SELECT id, city, filters FROM subscribe"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	subscriptions := []domain.Subscribe{}
	for rows.Next() {
		var subscription domain.Subscribe
		var rawFilters []byte
		err = rows.Scan(&subscription.ID, &subscription.City, &rawFilters)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(rawFilters, &subscription.Filters)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}
