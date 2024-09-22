package userstate

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func NewRedisClient() (*redis.Client, error) {
	log.Println("Opening Redis connection...")
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес Redis-сервера
		Password: "",               // Пароль (если не требуется, оставить пустым)
		DB:       0,                // Используемая база данных
	})

	log.Println("Pinging Redis...")
	_, err := client.Ping(ctx).Result() // Проверка соединения с Redis
	if err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}
