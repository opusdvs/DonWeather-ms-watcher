package usecase

import (
	"context"
	"time"

	"github.com/opusdvs/DonWeather-ms-watcher/internal/domain"
)

type WeatherService struct {
	weatherRepository      domain.WeatherRepository
	subscriptionRepository domain.SubscriptionRepository
	eventRepository        domain.EventRepository
}

func NewWeatherService(weatherRepository domain.WeatherRepository, subscriptionRepository domain.SubscriptionRepository, eventRepository domain.EventRepository) *WeatherService {
	return &WeatherService{
		weatherRepository:      weatherRepository,
		subscriptionRepository: subscriptionRepository,
		eventRepository:        eventRepository,
	}
}

func (s *WeatherService) GetLastStateWeather(ctx context.Context, city string) (*domain.StateWeather, error) {
	return s.weatherRepository.GetLastStateWeather(ctx, city)
}

func (s *WeatherService) GetCurrentStateWeather(ctx context.Context, city string) (*domain.StateWeather, error) {
	return s.weatherRepository.GetCurrentStateWeather(ctx, city)
}

func (s *WeatherService) GetActiveSubscriptions(ctx context.Context) ([]domain.Subscribe, error) {
	return s.subscriptionRepository.GetActiveSubscriptions(ctx)
}

func (s *WeatherService) CreateEventWeather(ctx context.Context, event *domain.EventWeather) error {
	return s.eventRepository.CreateEventWeather(ctx, event)
}

func (s *WeatherService) ProcessSubscriptions(ctx context.Context) error {
	subscriptions, err := s.GetActiveSubscriptions(ctx)
	if err != nil {
		return err
	}
	for _, subscription := range subscriptions {
		lastState, err := s.weatherRepository.GetLastStateWeather(ctx, subscription.City)
		if err != nil {
			continue
		}

		currentState, err := s.weatherRepository.GetCurrentStateWeather(ctx, subscription.City)
		if err != nil {
			continue
		}

		events, err := s.CompareStateWeather(ctx, &subscription, lastState, currentState)
		if err != nil {
			continue
		}
		for _, event := range events {
			s.eventRepository.CreateEventWeather(ctx, event)
		}
	}
	return nil
}

// ОЧень много однообразного кода, нужно переделать
func (s *WeatherService) CompareStateWeather(ctx context.Context, subscription *domain.Subscribe, lastState *domain.StateWeather, currentState *domain.StateWeather) ([]*domain.EventWeather, error) {
	events := []*domain.EventWeather{}
	if subscription.Filters.Temperature {
		if currentState.Temperature < lastState.Temperature+domain.TempRiseThreshold {
			events = append(events, &domain.EventWeather{
				EventType: domain.EventTempDrop,
				OldValue:  lastState.Temperature,
				NewValue:  currentState.Temperature,
				CreatedAt: time.Now(),
			})
		}
	} else if currentState.Temperature > lastState.Temperature-domain.TempDropThreshold {
		events = append(events, &domain.EventWeather{
			EventType: domain.EventTempRise,
			OldValue:  lastState.Temperature,
			NewValue:  currentState.Temperature,
			CreatedAt: time.Now(),
		})
	}
	if subscription.Filters.Humidity {
		if currentState.Humidity < lastState.Humidity+domain.HumidityChangeThreshold {
			events = append(events, &domain.EventWeather{
				EventType: domain.EventHumidityDrop,
				OldValue:  lastState.Humidity,
				NewValue:  currentState.Humidity,
				CreatedAt: time.Now(),
			})
		}
	} else if currentState.Humidity > lastState.Humidity-domain.HumidityChangeThreshold {
		events = append(events, &domain.EventWeather{
			EventType: domain.EventHumidityRise,
			OldValue:  lastState.Humidity,
			NewValue:  currentState.Humidity,
			CreatedAt: time.Now(),
		})
	}
	if subscription.Filters.WindSpeed {
		if currentState.WindSpeed < lastState.WindSpeed+domain.StrongWindThreshold {
			events = append(events, &domain.EventWeather{
				EventType: domain.EventStrongWind,
				OldValue:  lastState.WindSpeed,
				NewValue:  currentState.WindSpeed,
				CreatedAt: time.Now(),
			})
		}
	} else if currentState.WindSpeed > lastState.WindSpeed-domain.StrongWindThreshold {
		events = append(events, &domain.EventWeather{
			EventType: domain.EventStrongWind,
			OldValue:  lastState.WindSpeed,
			NewValue:  currentState.WindSpeed,
			CreatedAt: time.Now(),
		})
	}
	if subscription.Filters.Pressure {
		if currentState.Pressure < lastState.Pressure+domain.PressureChangeThreshold {
			events = append(events, &domain.EventWeather{
				EventType: domain.EventPressureDrop,
				OldValue:  lastState.Pressure,
				NewValue:  currentState.Pressure,
				CreatedAt: time.Now(),
			})
		}
	} else if currentState.Pressure > lastState.Pressure-domain.PressureChangeThreshold {
		events = append(events, &domain.EventWeather{
			EventType: domain.EventPressureRise,
			OldValue:  lastState.Pressure,
			NewValue:  currentState.Pressure,
			CreatedAt: time.Now(),
		})
	}
	return events, nil
}
