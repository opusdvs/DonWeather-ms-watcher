package domain

import "time"

type EventType string

const (
	TempRiseThreshold       = 5.0
	TempDropThreshold       = 5.0
	StrongWindThreshold     = 15.0
	RainStartThreshold      = 70.0
	RainStopThreshold       = 30.0
	PressureChangeThreshold = 10.0
	HumidityChangeThreshold = 20.0
)

const (
	EventTempDrop     EventType = "TEMP_DROP"
	EventTempRise     EventType = "TEMP_RISE"
	EventRainStarted  EventType = "RAIN_STARTED"
	EventRainStopped  EventType = "RAIN_STOPPED"
	EventStrongWind   EventType = "STRONG_WIND"
	EventHumidityDrop EventType = "HUMIDITY_DROP"
	EventHumidityRise EventType = "HUMIDITY_RISE"
	EventPressureDrop EventType = "PRESSURE_DROP"
	EventPressureRise EventType = "PRESSURE_RISE"
)

type EventWeatherMessage struct {
	ID         string       `json:"id"`
	TelegramID string       `json:"telegram_id"`
	Data       EventWeather `json:"data"`
	Timestamp  time.Time    `json:"timestamp"`
	Type       string       `json:"type"`
	Source     string       `json:"source"`
	Version    string       `json:"version"`
}

type Subscribe struct {
	ID         string  `json:"id"`
	City       string  `json:"city"`
	Filters    Filters `json:"filters"`
	TelegramID string  `json:"telegram_id"`
}

type Filters struct {
	Humidity    bool `json:"humidity"`
	Temperature bool `json:"temperature"`
	WindSpeed   bool `json:"wind_speed"`
	Pressure    bool `json:"pressure"`
}

type StateWeather struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		LastUpdated string  `json:"last_updated"`
		TempC       float64 `json:"temp_c"`
		Humidity    float64 `json:"humidity"`
		PressureMb  float64 `json:"pressure_mb"`
		WindKph     float64 `json:"wind_kph"`
		Condition   struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EventWeather struct {
	City      string    `json:"city"`
	EventType EventType `json:"event_type"`
	OldValue  float64   `json:"old_value"`
	NewValue  float64   `json:"new_value"`
	CreatedAt time.Time `json:"created_at"`
}
