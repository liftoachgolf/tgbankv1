package models

type AccountTg struct {
	id       int64  `json:"id"`
	chat_id  int64  `json:"chat_id"`
	username string `json:"username"`
	balance  int64  `json:"balance"`
	currency string `json:"currency"`
}

type MessegeTg struct {
	chat_id    int64 `json:"chat_id"`
	message_id int64 `json:"message_id"`
}
