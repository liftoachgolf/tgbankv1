package models

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int              `json:"update_id"`
	GetMessage    *IncomingMessage `json:"message"`        // Сообщения
	CallbackQuery *CallbackQuery   `json:"callback_query"` // Callback-запросы
}

type IncomingMessage struct {
	Text      string `json:"text"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
	MessageId int    `json:"message_id"`
}

type CallbackQuery struct {
	ID      string           `json:"id"`
	From    From             `json:"from"`
	Message *IncomingMessage `json:"message"`
	Data    string           `json:"data"` // Данные, связанные с callback-запросом
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type Message struct {
	MessageId    int
	Text         string
	ChatId       int
	Username     string
	CallbackData string
}

type UpdateMessage struct {
	ChatId    int    `json:"chat_id"`
	Text      string `json:"text"`
	MessageId int    `json:"message_id"`
}
