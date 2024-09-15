package postgrestg

import (
	"context"
	"database/sql"
	"fmt"
)

type storeTgSQL struct {
	db *sql.DB
	*Repository
}

func NewStoreSQL(db *sql.DB) StoreTgSQL {
	return &storeTgSQL{
		db: db,
	}
}

func (store *storeTgSQL) ExecTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
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
