package domain

import "context"

type EventRepository interface {
	CreateEventWeather(ctx context.Context, subscription Subscribe, event EventWeather) error
}
type WeatherRepository interface {
	GetLastStateWeather(ctx context.Context, city string) (StateWeather, error)
	GetCurrentStateWeather(ctx context.Context, city string) (StateWeather, error)
	SaveStateWeather(ctx context.Context, city string, stateWeather StateWeather) error
}

type SubscriptionRepository interface {
	GetActiveSubscriptions(ctx context.Context) ([]Subscribe, error)
}
