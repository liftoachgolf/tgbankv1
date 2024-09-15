package postgrestg

import (
	"context"
	"database/sql"
	"tgBank/models"
)

type AccountTgSQL interface {
	CreateAccount(ctx context.Context, id int64, username string, balance int64, currency string) error
	GetAccount(ctx context.Context, chatID int64) (*models.AccountTg, error)
	IsExitsts(ctx context.Context, chatID int64) (bool, error)
	AddMoney(ctx context.Context, chatId1 int64, amount1 int64, chatId2 int64, amount2 int64) (account1 models.AccountTg, account2 models.AccountTg, err error)
	AddAccountBalance(ctx context.Context, arg models.AddAccountTgBalanceParams) (models.AccountTg, error)
}

type StoreTgSQL interface {
	ExecTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}

type MessageTgSQL interface {
	CreateMessage(msg models.Message) error
}

type EntryTgSQL interface {
	CreateEntry(ctx context.Context, arg models.CreateEntryTgParams) (models.EntryTg, error)
	GetEntry(ctx context.Context, id int64) (models.EntryTg, error)
	ListEntries(ctx context.Context, arg models.ListEntriesTgParams) ([]models.EntryTg, error)
}

type TransferTgSQL interface {
	CreateTransfer(ctx context.Context, arg models.CreateTransferTgParams) (models.TransferTg, error)
	GetTransfer(ctx context.Context, id int64) (models.TransferTg, error)
	ListTransfers(ctx context.Context, arg models.ListTransfersTgParams) ([]models.TransferTg, error)
}

type Repository struct {
	AccountTgSQL
	MessageTgSQL
	StoreTgSQL
	EntryTgSQL
	TransferTgSQL
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		AccountTgSQL:  NewAccountTGSQL(db),
		MessageTgSQL:  NewMessageTgSQL(db),
		StoreTgSQL:    NewStoreSQL(db),
		EntryTgSQL:    NewEntrySQL(db),
		TransferTgSQL: NewTransferSQL(db),
	}
}
