package processor

import (
	"fmt"
	"log"
	"tgBank/models"
	"tgBank/pkg/telegram"
	"time"
)

type Processor interface {
	Fetch(limit int) ([]models.Message, error)
	processEvent(u models.Update) models.Message
	HandleMessage(msg models.Message) error
}

type processor struct {
	tgClient telegram.Client
	offset   int
}

func NewProcessor(tgClient telegram.Client) Processor {
	return &processor{
		tgClient: tgClient,
		offset:   0,
	}
}

// HandleMessage обрабатывает входящие сообщения
func (p *processor) HandleMessage(msg models.Message) error {
	log.Print(msg.MessageId)
	log.Print(msg.Username)

	if msg.Text == "/start" {
		err := p.tgClient.SendMessage(msg.ChatId, "саламуалейкум брат")
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		return p.tgClient.UpdateMessage(models.UpdateMessage{
			Text:      "дима пидар",
			ChatId:    msg.ChatId,
			MessageId: msg.MessageId + 1,
		})
	}
	return p.tgClient.SendMessage(msg.ChatId, msg.Text)
}

func (p *processor) Fetch(limit int) ([]models.Message, error) {
	updates, err := p.tgClient.GetUpdates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get messages: %w", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]models.Message, 0, len(updates))
	for _, u := range updates {
		res = append(res, p.processEvent(u))
	}

	// Обновляем смещение
	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *processor) processEvent(u models.Update) models.Message {
	msg := models.Message{
		Text:      fetchText(u),
		ChatId:    u.GetMessage.Chat.ID,
		Username:  u.GetMessage.From.Username,
		MessageId: u.GetMessage.MessageId,
	}

	// // Сохранение сообщения в базу данных
	// err := p.msgRepo.SaveMessage(msg)
	// if err != nil {
	// 	fmt.Printf("error saving message: %v\n", err)
	// }

	return msg
}

// fetchText получает текст сообщения из обновления
func fetchText(u models.Update) string {
	if u.GetMessage == nil {
		return ""
	}
	return u.GetMessage.Text
}
