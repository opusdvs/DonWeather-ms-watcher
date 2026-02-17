package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/nats-io/nats.go"
	"github.com/opusdvs/DonWeather-ms-watcher/internal/delivery"
	"github.com/opusdvs/DonWeather-ms-watcher/internal/delivery/middleware"
	"github.com/opusdvs/DonWeather-ms-watcher/internal/repository"
	"github.com/opusdvs/DonWeather-ms-watcher/internal/usecase"
	"github.com/redis/go-redis/v9"
)

func main() {
	natsHost := os.Getenv("NATS_HOST")
	if natsHost == "" {
		log.Fatal("NATS_HOST is not set")
	}

	weatherApiURL := os.Getenv("WEATHER_API_URL")
	if weatherApiURL == "" {
		log.Fatal("WEATHER_API_URL is not set")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal("DB_HOST is not set")
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER is not set")
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD is not set")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME is not set")
	}
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		log.Fatal("REDIS_HOST is not set")
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		log.Fatal("REDIS_PORT is not set")
	}
	db, err := sql.Open("pgx", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbPingCtx, dbPingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbPingCancel()
	err = db.PingContext(dbPingCtx)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPingCancel()

	natsConn, err := nats.Connect(natsHost)
	if err != nil {
		log.Fatal(err)
	}
	defer natsConn.Close()
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:   0,
	})
	defer redisClient.Close()
	redisPingCtx, redisPingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer redisPingCancel()
	pong, err := redisClient.Ping(redisPingCtx).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pong)

	eventRepo, err := repository.NewEventRepository(natsConn)
	if err != nil {
		log.Fatal(err)
	}
	weatherRepo := repository.NewWeatherRepository(redisClient, weatherApiURL)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	weatherService := usecase.NewWeatherService(weatherRepo, subscriptionRepo, eventRepo)
	weatherDelivery := delivery.NewWeatherHTTPHandler(*weatherService)

	weatherMux := http.NewServeMux()
	weatherMux.HandleFunc("/api/v1/process-subscriptions", weatherDelivery.ProcessSubscriptions)
	handleMiddleware := middleware.MiddlewareChain(weatherMux, middleware.CorsMiddleware)

	mainMux := http.NewServeMux()
	mainMux.Handle("/", handleMiddleware)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mainMux,
	}
	server.ListenAndServe()
}
