package handler

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"tgBank/internal/processor"
	"tgBank/models"
	"time"
)

type Consumer interface {
	Process(event models.Message) error
	Start() error
	Handle(text string, chatId int, username string, messageId int, clback string) error
	handleEvents(events []models.Message) error
	processMessage(event models.Message) error
}
type consumer struct {
	processor processor.Processor
	batchSize int
	workers   int
}

func NewConsumer(processor processor.Processor, batchSize int) Consumer {
	return &consumer{
		processor: processor,
		batchSize: batchSize,
		workers:   100,
	}
}
func (c *consumer) Start() error {
	for {
		gotEvents, err := c.processor.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)
			continue
		}
	}
}

func (c *consumer) Handle(text string, chatId int, username string, messageId int, clback string) error {
	text = strings.TrimSpace(text)
	msg := models.Message{
		ChatId:       chatId,
		Text:         text,
		MessageId:    messageId,
		Username:     username,
		CallbackData: clback,
	}
	log.Printf("got new command '%s' from '%s'", text, username)
	return c.processor.HandleMessage(msg)
}

func (c *consumer) handleEvents(events []models.Message) error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, c.workers) // Семафор для ограничения количества горутин

	for _, event := range events {
		wg.Add(1)
		sem <- struct{}{}
		go func(e models.Message) {
			defer wg.Done()
			defer func() { <-sem }() // Освобождаем место в семафоре

			if err := c.Process(e); err != nil {
				log.Printf("can't handle event: %s", err.Error())
			}
		}(event)
	}

	wg.Wait() // Ожидаем завершения всех горутин
	return nil
}

func (c *consumer) processMessage(event models.Message) error {
	if err := c.Handle(event.Text, event.ChatId, event.Username, event.MessageId, event.CallbackData); err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}
	return nil
}

func (c *consumer) Process(event models.Message) error {
	return c.processMessage(event)
}
