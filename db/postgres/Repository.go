package postgres

import (
	"context"
	"database/sql"
)

type AccountSQL interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
	addMoney(ctx context.Context, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 Account, account2 Account, err error)
}

type EntrySQL interface {
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	GetEntry(ctx context.Context, id int64) (Entry, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
}

type TransferSQL interface {
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	GetTransfer(ctx context.Context, id int64) (Transfer, error)
}

type Repository struct {
	AccountSQL
	EntrySQL
	TransferSQL
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		AccountSQL:  NewAccountDb(db),
		EntrySQL:    NewEntryDb(db),
		TransferSQL: NewTransferDb(db),
	}
}
