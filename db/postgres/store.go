package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"tgBank/models"
)

type StoreSQLImpl struct {
	db *sql.DB
	*Repository
}

func NewStoreSQL(db *sql.DB) *StoreSQLImpl {
	return &StoreSQLImpl{
		db: db,
	}
}

func (store *StoreSQLImpl) execTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

func (store *StoreSQLImpl) TransferTx(ctx context.Context, arg models.TransferTxParams) (models.TransferTxResult, error) {
	var result models.TransferTxResult

	err := store.execTx(ctx, func(tx *sql.Tx) error {
		var err error

		result.Transfer, err = store.CreateTransfer(ctx, models.CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = store.CreateEntry(ctx, models.CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = store.CreateEntry(ctx, models.CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = store.addMoney(ctx, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = store.addMoney(ctx, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

		}

		return nil
	})

	return result, err
}
