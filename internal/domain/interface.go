package domain

import "context"

type EventRepository interface {
	CreateEventWeather(ctx context.Context, event *EventWeather) error
}
type WeatherRepository interface {
	GetLastStateWeather(ctx context.Context, city string) (*StateWeather, error)
	GetCurrentStateWeather(ctx context.Context, city string) (*StateWeather, error)
}

type SubscriptionRepository interface {
	GetActiveSubscriptions(ctx context.Context) ([]Subscribe, error)
}
