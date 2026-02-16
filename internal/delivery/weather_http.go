package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/opusdvs/DonWeather-ms-watcher/internal/usecase"
)

type WeatherHTTPHandler struct {
	weatherService usecase.WeatherService
}

func NewWeatherHTTPHandler(weatherService usecase.WeatherService) *WeatherHTTPHandler {
	return &WeatherHTTPHandler{weatherService: weatherService}
}

func (h *WeatherHTTPHandler) GetLastStateWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	state, err := h.weatherService.GetLastStateWeather(r.Context(), city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(state)
}
