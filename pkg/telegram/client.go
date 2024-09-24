package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"tgBank/models"
)

type Client interface {
	SendMessage(chatID int, text string) error
	GetUpdates(offset int, limit int) ([]models.Update, error)
	SetChatMenuButton(chatId int, button models.MenuButton) error
	UpdateMessage(updateMsg models.UpdateMessage) error
	SetCommands(commands models.Commands) error
	SetInlineButton(chatID int, keyboard models.InlineKeyboardMarkup, text string) error
}

type client struct {
	host     string
	basePath string
	client   http.Client
}

func NewTelegramApiClient(host, token string) Client {
	return &client{
		host:     host,
		basePath: "bot" + token,
		client:   http.Client{},
	}
}

// SendMessage отправляет сообщение в Telegram чат
func (r *client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := r.doRequest("sendMessage", q)
	return err
}

// GetUpdates получает обновления (например, сообщения) от Telegram API
func (r *client) GetUpdates(offset int, limit int) ([]models.Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := r.doRequest("getUpdates", q)
	if err != nil {
		return nil, err
	}

	var res models.UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

// SetChatMenuButton устанавливает кнопку меню для чата
func (r *client) SetChatMenuButton(chatId int, button models.MenuButton) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	buttonJson, err := json.Marshal(button)
	if err != nil {
		return fmt.Errorf("error marshaling menu button: %v", err)
	}
	q.Add("menu_button", string(buttonJson))
	_, err = r.doRequest("setChatMenuButton", q)
	return err
}

// UpdateMessage обновляет текст сообщения
func (r *client) UpdateMessage(updateMsg models.UpdateMessage) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(updateMsg.ChatId))
	q.Add("message_id", strconv.Itoa(updateMsg.MessageId))
	q.Add("text", updateMsg.Text)

	_, err := r.doRequest("editMessageText", q)
	return err
}

func (r *client) SetInlineButton(chatID int, keyboard models.InlineKeyboardMarkup, text string) error {
	keyboardJSON, err := json.Marshal(keyboard)
	if err != nil {
		return fmt.Errorf("error marshalling keyboard: %w", err)
	}

	params := url.Values{}
	params.Add("chat_id", fmt.Sprintf("%d", chatID))
	params.Add("text", text)
	params.Add("reply_markup", string(keyboardJSON))
	_, err = r.doRequest("sendMessage", params)
	return err
}

// SetCommands устанавливает команды для бота
func (r *client) SetCommands(commands models.Commands) error {
	q := url.Values{}
	commandsJson, err := json.Marshal(commands)
	if err != nil {
		return fmt.Errorf("error while setting commands: %s", err)
	}
	q.Add("commands", string(commandsJson))
	_, err = r.doRequest("setMyCommands", q)
	return err
}

// doRequest - вспомогательный метод для отправки запросов к Telegram API
func (r *client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   r.host,
		Path:   path.Join(r.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("cant do request: %w", err)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cant do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't do request: received non-200 status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cant do request: %w", err)
	}
	return body, nil
}
