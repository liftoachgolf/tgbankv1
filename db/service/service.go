package service

import (
	"context"
	"database/sql"
	"tgBank/db/postgres"
	"tgBank/models"
)

type AccountService interface {
	CreateAccount(ctx context.Context, arg models.CreateAccountParams) (models.Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (models.Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (models.Account, error)
	ListAccounts(ctx context.Context, arg models.ListAccountsParams) ([]models.Account, error)
	UpdateAccount(ctx context.Context, arg models.UpdateAccountParams) (models.Account, error)
	addMoney(ctx context.Context, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 models.Account, account2 models.Account, err error)
}

type EntryService interface {
	CreateEntry(ctx context.Context, arg models.CreateEntryParams) (models.Entry, error)
	GetEntry(ctx context.Context, id int64) (models.Entry, error)
	ListEntries(ctx context.Context, arg models.ListEntriesParams) ([]models.Entry, error)
}

type TransferService interface {
	ListTransfers(ctx context.Context, arg models.ListTransfersParams) ([]models.Transfer, error)
	CreateTransfer(ctx context.Context, arg models.CreateTransferParams) (models.Transfer, error)
	GetTransfer(ctx context.Context, id int64) (models.Transfer, error)
}

type StoreService interface {
	TransferTx(ctx context.Context, arg models.TransferTxParams) (models.TransferTxResult, error)
	execTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}

type Service struct {
	AccountService
	EntryService
	TransferService
}

func NewRepository(repo *postgres.Repository) *Service {
	return &Service{
		AccountService:  NewAccountServiceImpl(repo.AccountSQL),
		EntryService:    NewEntryServiceImpl(repo.EntrySQL),
		TransferService: NewTransferServiceImpl(repo.TransferSQL),
	}
}
