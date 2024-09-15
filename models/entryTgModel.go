package models

import "time"

type EntryTg struct {
	ID        int64     `json:"id"`
	ChatId    int64     `json:"chat_id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateEntryTgParams struct {
	ChatId int64 `json:"chat_id"`
	Amount int64 `json:"amount"`
}
type ListEntriesTgParams struct {
	ChatId int64 `json:"chat_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}
