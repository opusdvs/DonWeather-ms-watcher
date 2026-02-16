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

type Subscribe struct {
	ID      string  `json:"id"`
	City    string  `json:"city"`
	Filters Filters `json:"filters"`
}

type Filters struct {
	Humidity    bool `json:"humidity"`
	Temperature bool `json:"temperature"`
	WindSpeed   bool `json:"wind_speed"`
	Pressure    bool `json:"pressure"`
}

type StateWeather struct {
	City        string    `json:"city"`
	Temperature float64   `json:"temperature"`
	WindSpeed   float64   `json:"wind_speed"`
	Humidity    float64   `json:"humidity"`
	Pressure    float64   `json:"pressure"`
	Clouds      float64   `json:"clouds"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EventWeather struct {
	EventType EventType `json:"event_type"`
	OldValue  float64   `json:"old_value"`
	NewValue  float64   `json:"new_value"`
	CreatedAt time.Time `json:"created_at"`
}
