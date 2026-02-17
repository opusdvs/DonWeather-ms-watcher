package delivery

import (
	"log"
	"net/http"

	"github.com/opusdvs/DonWeather-ms-watcher/internal/usecase"
)

type WeatherHTTPHandler struct {
	weatherService usecase.WeatherService
}

func NewWeatherHTTPHandler(weatherService usecase.WeatherService) *WeatherHTTPHandler {
	return &WeatherHTTPHandler{weatherService: weatherService}
}

func (h *WeatherHTTPHandler) ProcessSubscriptions(w http.ResponseWriter, r *http.Request) {
	err := h.weatherService.ProcessSubscriptions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error processing subscriptions:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Subscriptions processed successfully"))
}
