package models

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID         int              `json:"update_id"`
	GetMessage *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text      string `json:"text"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
	MessageId int    `json:"message_id"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type Message struct {
	MessageId int
	Text      string
	ChatId    int
	Username  string
}

type UpdateMessage struct {
	ChatId    int    `json:"chat_id"`
	Text      string `json:"text"`
	MessageId int    `json:"message_id"`
}
