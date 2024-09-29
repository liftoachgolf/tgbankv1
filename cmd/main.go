package main

import (
	"fmt"
	"log"
	"tgBank/internal/handler"
	processstateuser "tgBank/internal/processStateUser"
	"tgBank/internal/processor"
	serviceTg "tgBank/internal/servivceTg"
	postgrestg "tgBank/pkg/postgresTg"
	"tgBank/pkg/telegram"
	userstate "tgBank/pkg/userStateRepo"

	_ "github.com/lib/pq"
)

func main() {
	t := "xxx:xxx"
	log.Print("starting...")

	// Подключение к базе данных
	db, err := postgrestg.NewPostgresDb()
	if err != nil {
		log.Fatalf("error while connecting to db: %v", err)
	}
	dbR, err := userstate.NewRedisClient()
	if err != nil {
		log.Fatalf("error while connecting to db: %v", err)
	}

	// Создание репозиториев
	dbRepo := postgrestg.NewRepository(db)
	telegramRepo := telegram.NewTelegramApiClient("api.telegram.org", t)
	userStRepo := userstate.NewUserStateRepository(dbR)

	// Создание сервисов
	dbSrv := serviceTg.NewServiceTg(dbRepo)
	userStSrv := processstateuser.NewProcessUserStateService(userStRepo)
	telegramSvc := processor.NewProcessor(telegramRepo, dbSrv, userStSrv)

	// Создание хендлера Telegram
	telegramHandler := handler.NewConsumer(telegramSvc, 100)

	// Обработка обновлений
	err = telegramHandler.Start()
	if err != nil {
		log.Fatalf("failed to start handler: %v", err)
	}

	fmt.Print("started")
}
