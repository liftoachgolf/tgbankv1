package main

import (
	"fmt"
	"log"
	"tgBank/internal/handler"
	"tgBank/internal/processor"
	serviceTg "tgBank/internal/servivceTg"
	postgrestg "tgBank/pkg/postgresTg"
	"tgBank/pkg/telegram"

	_ "github.com/lib/pq"
)

func main() {
	t := "7077450094:AAExNrbVTytTFuEsOvb_aCldN5jk6RTyrAk"
	log.Print("starting...")

	// Подключение к базе данных
	db, err := postgrestg.NewPostgresDb()
	if err != nil {
		log.Fatalf("error while connecting to db: %v", err)
	}

	// Создание репозиториев
	dbRepo := postgrestg.NewRepository(db)
	telegramRepo := telegram.NewTelegramApiClient("api.telegram.org", t)

	// Создание сервисов
	dbSrv := serviceTg.NewServiceTg(dbRepo)
	telegramSvc := processor.NewProcessor(telegramRepo, dbSrv)

	// Создание хендлера Telegram
	telegramHandler := handler.NewConsumer(telegramSvc, 100)

	// Обработка обновлений
	err = telegramHandler.Start()
	if err != nil {
		log.Fatalf("failed to start handler: %v", err)
	}

	fmt.Print("started")
}
