package main

import (
	"log"
	"tgBank/internal/handler"
	"tgBank/internal/processor"
	"tgBank/pkg/telegram"

	_ "github.com/lib/pq"
)

func main() {
	t := "7077450094:AAExNrbVTytTFuEsOvb_aCldN5jk6RTyrAk"
	log.Print("starting...")

	// db, err := repository.NewPostgresDb()
	// if err != nil {
	// 	fmt.Errorf("error while connecting db: %w", err)
	// }

	// rdb, err := repository.NewRedisClient()
	// if err != nil {
	// 	fmt.Errorf("error initialize redis db: %w", err)
	// }
	// Создание репозитория Telegram
	telegramRepo := telegram.NewTelegramApiClient("api.telegram.org", t)

	// Создание сервиса Telegram
	telegramSvc := processor.NewProcessor(telegramRepo)

	// Создание хендлера Telegram
	telegramHandler := handler.NewConsumer(telegramSvc, 100)

	// Обработка обновлений
	err := telegramHandler.Start()
	if err != nil {
		panic(err)
	}
}
