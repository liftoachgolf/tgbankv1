package models

import "time"

type TransferTg struct {
	ID         int64     `json:"id"`
	FromChatId int64     `json:"from_chat_id"`
	ToChatId   int64     `json:"to_chat_id"`
	Amount     int64     `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateTransferTgParams struct {
	FromChatId int64 `json:"from_chat_id"`
	ToChatId   int64 `json:"to_chat_id"`
	Amount     int64 `json:"amount"`
}

type ListTransfersTgParams struct {
	FromChatId int64 `json:"from_chat_id"`
	ToChatId   int64 `json:"to_chat_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}
