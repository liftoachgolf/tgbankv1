package models

type TransferTgTxParams struct {
	FromChatId int64 `json:"from_chat_id"`
	ToChatId   int64 `json:"to_chat_id"`
	Amount     int64 `json:"amount"`
}

type TransferTgTxResult struct {
	Transfer   TransferTg `json:"transfer"`
	FromChatId AccountTg  `json:"from_chat_id"`
	ToChatId   AccountTg  `json:"to_chat_id"`
	FromEntry  EntryTg    `json:"from_entry"`
	ToEntry    EntryTg    `json:"to_entry"`
}
