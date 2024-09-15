package models

type AccountTg struct {
	ID       int64  `json:"id"`
	ChatID   int64  `json:"chat_id"`
	Username string `json:"username"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type AddAccountTgBalanceParams struct {
	Amount int64 `json:"amount"`
	ChatID int64 `json:"chat_id"`
}
