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

func (p *processor) showActionMenu(chatID int) error {
	buttons := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Пополнить", CallbackData: "deposit"},
				{Text: "Перевести", CallbackData: "transfer"},
			},
		},
	}
	err := p.tgClient.SetInlineButton(chatID, buttons, "Выберите действие:")
	if err != nil {
		return fmt.Errorf("error sending action menu: %w", err)
	}
	return nil
}

func (p *processor) HandleMessage(msg models.Message) error {
	log.Print(msg.MessageId)
	log.Print(msg.Username)
	fmt.Printf("chat id: %d\n", msg.ChatId)

	if strings.ToLower(msg.Text) == "back" {
		err := p.tgUserState.SetState(int64(msg.ChatId), userstate.UserState{State: "START"})
		if err != nil {
			return fmt.Errorf("error resetting user state: %w", err)
		}
		return p.tgClient.SendMessage(msg.ChatId, "Возвращаемся в стартовое меню.")
	}

	msgs := strings.Fields(msg.Text)
	command := msgs[0]

	if command == "Начать" {
		err := p.tgUserState.SetState(int64(msg.ChatId), userstate.UserState{State: "WAITING_FOR_HELLO"})
		if err != nil {
			return fmt.Errorf("error while setting user state")
		}
	}

	if command == "Действия" {
		err := p.tgUserState.SetState(int64(msg.ChatId), userstate.UserState{State: "WaitingForTx"})
		if err != nil {
			return fmt.Errorf("error while setting user state")
		}
		// Отправка сообщения с inline-клавиатурой
		keyboard := models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Пополнить", CallbackData: "popolnit"},
					{Text: "Перевести", CallbackData: "perevesti"},
				},
				{
					{Text: "Назад", CallbackData: "back"},
				},
			},
		}
		return p.tgClient.SetInlineButton(msg.ChatId, keyboard, "Выберите действие:")
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

	case "WaitingForTx":
		err := p.tgClient.SendMessage(msg.ChatId, "Введите сумму для пополнения\nНапжмите Back, если хотите вернутся в стартовое меню")
		if err != nil {
			return fmt.Errorf("error while send msg: %s", err)
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

	// Обработка команд, которые пришли из inline-клавиатуры
	switch msg.CallbackData {
	case "popolnit":
		err := p.tgClient.SendMessage(msg.ChatId, "Введите сумму для пополнения:")
		if err != nil {
			return err
		}
		return p.tgUserState.SetState(int64(msg.ChatId), userstate.UserState{State: "POPOLNIT"})
	case "perevesti":
		err := p.tgClient.SendMessage(msg.ChatId, "Введите ID получателя и сумму для перевода:")
		if err != nil {
			return err
		}
		return p.tgUserState.SetState(int64(msg.ChatId), userstate.UserState{State: "PEREVESTI"})
	case "back":
		err := p.tgUserState.SetState(int64(msg.ChatId), userstate.UserState{State: "START"})
		if err != nil {
			return fmt.Errorf("error resetting user state: %w", err)
		}
		return p.tgClient.SendMessage(msg.ChatId, "Возвращаемся в стартовое меню.")
	}

	if command == "Пополнить" {
		amount, err := strconv.ParseInt(msgs[1], 10, 64)
		if err != nil {
			p.tgUserState.SetState(int64(msg.ChatId), userstate.UserState{State: "START"})
			return p.tgClient.SendMessage(msg.ChatId, fmt.Sprintf("Не является числом: %s ", msgs[1]))
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
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]models.Message, 0, len(updates))
	for _, u := range updates {
		msg := p.processEvent(u)
		if msg.ChatId != 0 {
			res = append(res, msg)
		}
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}
func (p *processor) HandleCallbackQuery(callback models.CallbackQuery) error {
	switch callback.Data {
	case "deposit":
		err := p.tgClient.SendMessage(callback.Message.Chat.ID, "Вы выбрали Пополнить. Введите сумму:")
		if err != nil {
			return fmt.Errorf("error sending deposit message: %w", err)
		}
		err = p.tgUserState.SetState(int64(callback.Message.Chat.ID), userstate.UserState{State: "WAITING_FOR_DEPOSIT_AMOUNT"})
		if err != nil {
			return fmt.Errorf("error setting user state: %w", err)
		}
	case "transfer":
		err := p.tgClient.SendMessage(callback.Message.Chat.ID, "Вы выбрали Перевести. Введите сумму и получателя:")
		if err != nil {
			return fmt.Errorf("error sending transfer message: %w", err)
		}
		err = p.tgUserState.SetState(int64(callback.Message.Chat.ID), userstate.UserState{State: "WAITING_FOR_TRANSFER_DETAILS"})
		if err != nil {
			return fmt.Errorf("error setting user state: %w", err)
		}
	default:
		err := p.tgClient.SendMessage(callback.Message.Chat.ID, "Неизвестное действие.")
		if err != nil {
			return fmt.Errorf("error handling callback: %w", err)
		}
	}
	return nil
}

func (p *processor) processEvent(u models.Update) models.Message {
	if u.CallbackQuery != nil {
		err := p.HandleCallbackQuery(*u.CallbackQuery)
		if err != nil {
			log.Printf("Error handling callback query: %v", err)
		}
		return models.Message{} // возвращаем пустое сообщение, так как это callback
	}

	if u.GetMessage == nil {
		return models.Message{}
	}

	msg := models.Message{
		Text:      fetchText(u),
		ChatId:    u.GetMessage.Chat.ID,
		Username:  u.GetMessage.From.Username,
		MessageId: u.GetMessage.MessageId,
	}
	return msg
}

func fetchText(u models.Update) string {
	if u.GetMessage != nil {
		return u.GetMessage.Text
	}
	return ""
}
