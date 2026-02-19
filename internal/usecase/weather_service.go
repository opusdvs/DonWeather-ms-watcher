package usecase

import (
	"context"
	"log"
	"math/rand/v2"
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

func (s *WeatherService) GetLastStateWeather(ctx context.Context, city string) (domain.StateWeather, error) {
	return s.weatherRepository.GetLastStateWeather(ctx, city)
}

func (s *WeatherService) GetCurrentStateWeather(ctx context.Context, city string) (domain.StateWeather, error) {
	return s.weatherRepository.GetCurrentStateWeather(ctx, city)
}

func (s *WeatherService) GetActiveSubscriptions(ctx context.Context) ([]domain.Subscribe, error) {
	return s.subscriptionRepository.GetActiveSubscriptions(ctx)
}

func (s *WeatherService) CreateEventWeather(ctx context.Context, subscription domain.Subscribe, event domain.EventWeather) error {
	return s.eventRepository.CreateEventWeather(ctx, subscription, event)
}

func (s *WeatherService) SaveStateWeather(ctx context.Context, city string, stateWeather domain.StateWeather) error {
	return s.weatherRepository.SaveStateWeather(ctx, city, stateWeather)
}

func (s *WeatherService) ProcessSubscriptions(ctx context.Context) error {
	log.Println("Processing subscriptions")
	subscriptions, err := s.GetActiveSubscriptions(ctx)
	if err != nil {
		log.Println("Error getting active subscriptions:", err)
		return err
	}
	for _, subscription := range subscriptions {
		log.Println("Processing subscription:", subscription.City)
		lastState, err := s.GetLastStateWeather(ctx, subscription.City)
		if err != nil {
			log.Println("Error getting last state weather (will treat as first run):", err)
		}
		currentState, err := s.GetCurrentStateWeather(ctx, subscription.City)
		if err != nil {
			log.Println("Error getting current state weather:", err)
			return err
		}
		if lastState.Location.Name == "" {
			log.Println("No last state weather for subscription:", subscription.City)
			err = s.SaveStateWeather(ctx, subscription.City, currentState)
			if err != nil {
				log.Println("Error saving state weather:", err)
				return err
			}
			continue
		}
		events, err := s.CompareStateWeather(ctx, subscription, &lastState, currentState)
		if err != nil {
			log.Println("Error comparing state weather:", err)
			return err
		}
		for _, event := range events {
			err = s.eventRepository.CreateEventWeather(ctx, subscription, event)
			if err != nil {
				log.Println("Error creating event weather:", err)
				return err
			}
		}
		err = s.SaveStateWeather(ctx, subscription.City, currentState)
		if err != nil {
			log.Println("Error saving state weather:", err)
			return err
		}
	}
	return nil
}

// Очень много однообразного кода, нужно переделать
func (s *WeatherService) CompareStateWeather(
	ctx context.Context,
	subscription domain.Subscribe,
	lastState *domain.StateWeather,
	currentState domain.StateWeather,
) ([]domain.EventWeather, error) {

	// Если нет предыдущего состояния — нечего сравнивать
	if lastState == nil {
		return nil, nil
	}

	var events []domain.EventWeather
	now := time.Now()

	// Универсальная функция для проверки изменений
	check := func(enabled bool, oldVal, newVal, riseThreshold, dropThreshold float64, riseType, dropType domain.EventType) {
		if !enabled {
			return
		}
		if newVal > oldVal+riseThreshold {
			events = append(events, domain.EventWeather{
				City:      subscription.City,
				EventType: riseType,
				OldValue:  oldVal,
				NewValue:  newVal,
				CreatedAt: now,
			})
		} else if newVal < oldVal-dropThreshold {
			events = append(events, domain.EventWeather{
				City:      subscription.City,
				EventType: dropType,
				OldValue:  oldVal,
				NewValue:  newVal,
				CreatedAt: now,
			})
		}
	}

	// Проверяем температуру
	check(subscription.Filters.Temperature,
		lastState.Current.TempC,
		currentState.Current.TempC,
		domain.TempRiseThreshold,
		domain.TempDropThreshold,
		domain.EventTempRise,
		domain.EventTempDrop,
	)

	// Проверяем влажность
	check(subscription.Filters.Humidity,
		lastState.Current.Humidity,
		currentState.Current.Humidity,
		domain.HumidityChangeThreshold,
		domain.HumidityChangeThreshold,
		domain.EventHumidityRise,
		domain.EventHumidityDrop,
	)

	// Проверяем скорость ветра
	check(subscription.Filters.WindSpeed,
		lastState.Current.WindKph,
		currentState.Current.WindKph,
		domain.StrongWindThreshold,
		domain.StrongWindThreshold,
		domain.EventStrongWind,
		domain.EventStrongWind, // можно разделить на Rise/Drop, если нужно
	)

	// Проверяем давление
	check(subscription.Filters.Pressure,
		lastState.Current.PressureMb,
		currentState.Current.PressureMb,
		domain.PressureChangeThreshold,
		domain.PressureChangeThreshold,
		domain.EventPressureRise,
		domain.EventPressureDrop,
	)

	return events, nil
}

func (s *WeatherService) CompareStateWeatherTest(
	subscription domain.Subscribe,
) []domain.EventWeather {
	log.Println("Comparing state weather test for subscription:", subscription.City)
	now := time.Now()
	events := []domain.EventWeather{}

	type checkDef struct {
		enabled bool
		name    string
		rise    domain.EventType
		drop    domain.EventType
	}

	// Собираем фильтры
	checks := []checkDef{
		{subscription.Filters.Temperature, "Temperature", domain.EventTempRise, domain.EventTempDrop},
		{subscription.Filters.Humidity, "Humidity", domain.EventHumidityRise, domain.EventHumidityDrop},
		{subscription.Filters.WindSpeed, "WindSpeed", domain.EventStrongWind, domain.EventStrongWind},
		{subscription.Filters.Pressure, "Pressure", domain.EventPressureRise, domain.EventPressureDrop},
	}

	for _, c := range checks {
		if !c.enabled {
			continue
		}

		// Генерируем случайное значение для теста
		newVal := rand.Float64() * 100
		oldVal := rand.Float64() * 100

		eventType := c.rise
		if newVal < oldVal {
			eventType = c.drop
		}

		events = append(events, domain.EventWeather{
			City:      subscription.City,
			EventType: eventType,
			OldValue:  oldVal,
			NewValue:  newVal,
			CreatedAt: now,
		})
	}

	return events
}
