package processor

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	serviceTg "tgBank/internal/servivceTg"
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
	tgClient  telegram.Client
	tgService *serviceTg.ServiceTg // Добавляем ссылку на сервис работы с базой данных
	offset    int
}

func NewProcessor(tgClient telegram.Client, tgService *serviceTg.ServiceTg) Processor {
	return &processor{
		tgClient:  tgClient,
		tgService: tgService, // Инициализируем сервис
		offset:    0,
	}
}

// HandleMessage обрабатывает входящие сообщения
func (p *processor) HandleMessage(msg models.Message) error {
	log.Print(msg.MessageId)
	log.Print(msg.Username)
	fmt.Printf("chat id: %d\n", msg.ChatId)

	msgs := strings.Fields(msg.Text)
	if msg.Text == "/start" {
		err := p.tgService.MessageTgService.CreateMessage(msg)
		if err != nil {
			return fmt.Errorf("error while adding msg to db: %w ", err)
		}
		err = p.tgClient.SendMessage(msg.ChatId, "саламуалейкум брат, ваш аккаунт создан")
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

	if msgs[0] == "Пополнить" {
		amount, err := strconv.ParseInt(msgs[1], 10, 64)
		if err != nil {
			return err
		}
		acc, err := p.tgService.AccountTgService.AddAccountBalance(context.Background(), models.AddAccountTgBalanceParams{
			Amount: amount,
			ChatID: int64(msg.ChatId),
		})
		if err != nil {
			return err
		}
		return p.tgClient.SendMessage(msg.ChatId, fmt.Sprint("ur new balance: ", acc.Balance))
	}
	if msgs[0] == "Перевести" {
		toChat, err := strconv.ParseInt(msgs[1], 10, 64)
		if err != nil {
			return err
		}
		amount, err := strconv.ParseInt(msgs[2], 10, 64)
		if err != nil {
			return err
		}
		res, err := p.tgService.StoreTgService.TransferTx(context.Background(), models.TransferTgTxParams{
			FromChatId: int64(msg.ChatId),
			ToChatId:   toChat,
			Amount:     amount,
		})
		fmt.Print(res)
		return p.tgClient.SendMessage(msg.ChatId, fmt.Sprintf("теперь твой баланс равен: %d, а баланс получателя: %d", res.FromChatId.Balance, res.ToChatId.Balance))
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
	return msg
}

// fetchText получает текст сообщения из обновления
func fetchText(u models.Update) string {
	if u.GetMessage == nil {
		return ""
	}
	return u.GetMessage.Text
}
