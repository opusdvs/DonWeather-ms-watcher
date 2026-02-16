package repository

import (
	"context"
	"database/sql"

	"github.com/opusdvs/DonWeather-ms-watcher/internal/domain"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) GetActiveSubscriptions(ctx context.Context) ([]domain.Subscribe, error) {
	return nil, nil
}
