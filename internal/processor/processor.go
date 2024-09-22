package processor

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	processstateuser "tgBank/internal/processStateUser"
	serviceTg "tgBank/internal/servivceTg"
	"tgBank/models"
	"tgBank/pkg/telegram"
	userstate "tgBank/pkg/userStateRepo"
	"time"
)

type Processor interface {
	Fetch(limit int) ([]models.Message, error)
	processEvent(u models.Update) models.Message
	HandleMessage(msg models.Message) error
}

type processor struct {
	tgClient    telegram.Client
	tgService   *serviceTg.ServiceTg
	tgUserState *processstateuser.ProcessUserStateService
	offset      int
}

func NewProcessor(tgClient telegram.Client, tgService *serviceTg.ServiceTg, tgUserState *processstateuser.ProcessUserStateService) Processor {
	return &processor{
		tgClient:    tgClient,
		tgService:   tgService, // Инициализируем сервис
		tgUserState: tgUserState,
		offset:      0,
	}
}

func (p *processor) HandleMessage(msg models.Message) error {
	log.Print(msg.MessageId)
	log.Print(msg.Username)
	fmt.Printf("chat id: %d\n", msg.ChatId)

	msgs := strings.Fields(msg.Text)
	command := msgs[0]

	if command == "Начать" {
		err := p.tgUserState.SetState(int64(msg.ChatId), userstate.UserState{State: "WAITING_FOR_HELLO"})
		if err != nil {
			fmt.Errorf("error while setting user state")
		}
	}

	state, err := p.tgUserState.UserStateService.GetState(int64(msg.ChatId))
	if err != nil {
		log.Printf("Ошибка получения состояния пользователя: %v", err)

		state = userstate.UserState{State: "START"}
	}

	switch state.State {
	case "START":
		if command == "/start" {

			err := p.tgService.MessageTgService.CreateMessage(msg)
			if err != nil {
				return fmt.Errorf("error while adding msg to db: %w ", err)
			}
			err = p.tgClient.SendMessage(msg.ChatId, "ваш аккаунт зарегистрирован")
			if err != nil {
				return err
			}
			time.Sleep(1 * time.Second)
			err = p.tgClient.UpdateMessage(models.UpdateMessage{
				Text:      "Привет! Я ваш бот.",
				ChatId:    msg.ChatId,
				MessageId: msg.MessageId + 1,
			})
			if err != nil {
				return err
			}

			err = p.tgUserState.SetState(int64(msg.ChatId), userstate.UserState{State: "WAITING_FOR_HELLO"})
			if err != nil {
				return fmt.Errorf("error setting user state: %w", err)
			}
			return nil
		}

	case "WAITING_FOR_HELLO":
		if strings.ToLower(msg.Text) == "привет" {
			err := p.tgClient.SendMessage(msg.ChatId, "Как дела?")
			if err != nil {
				return err
			}
			err = p.tgUserState.SetState(int64(msg.ChatId), userstate.UserState{State: "END"})
			if err != nil {
				return fmt.Errorf("error setting user state: %w", err)
			}
			return nil
		} else {
			err := p.tgClient.SendMessage(msg.ChatId, "Пожалуйста, напишите 'привет'.")
			if err != nil {
				return err
			}
			return nil
		}

	case "END":
		err := p.tgClient.SendMessage(msg.ChatId, "До свидания!")
		if err != nil {
			return err
		}
		// Сброс состояния
		err = p.tgUserState.SetState(int64(msg.ChatId), userstate.UserState{State: "START"})
		if err != nil {
			return fmt.Errorf("error resetting user state: %w", err)
		}
		return nil

	default:

		err := p.tgClient.SendMessage(msg.ChatId, "Неизвестное состояние. Сброс состояния.")
		if err != nil {
			return err
		}
		// Сброс состояния
		err = p.tgUserState.SetState(int64(msg.ChatId), userstate.UserState{State: "START"})
		if err != nil {
			return fmt.Errorf("error resetting user state: %w", err)
		}
		return nil
	}

	if command == "Пополнить" {
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
		return p.tgClient.SendMessage(msg.ChatId, fmt.Sprintf("Ваш новый баланс: %d", acc.Balance))
	}
	if command == "Перевести" {
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
		if err != nil {
			return err
		}
		return p.tgClient.SendMessage(msg.ChatId, fmt.Sprintf("Теперь ваш баланс равен: %d, а баланс получателя: %d", res.FromChatId.Balance, res.ToChatId.Balance))
	}

	// Обработка неизвестных команд
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

func fetchText(u models.Update) string {
	if u.GetMessage == nil {
		return ""
	}
	return u.GetMessage.Text
}
