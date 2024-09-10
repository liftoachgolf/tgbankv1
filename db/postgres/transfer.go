package postgres

import (
	"context"
	"database/sql"
)

type TransferDb struct {
	db *sql.DB
}

func NewTransferDb(db *sql.DB) *TransferDb {
	return &TransferDb{
		db: db,
	}
}

type CreateTransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func (q *TransferDb) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	query := `INSERT INTO transfers(
    from_account_id,
    to_account_id,
    amount
)VALUES(
    $1, $2, $3
) RETURNING id, from_account_id, to_account_id, amount, created_at`
	row := q.db.QueryRowContext(ctx, query, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

func (q *TransferDb) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	query := `SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE id = $1
LIMIT 1`
	row := q.db.QueryRowContext(ctx, query, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

type ListTransfersParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Limit         int32 `json:"limit"`
	Offset        int32 `json:"offset"`
}

func (q *TransferDb) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error) {
	query := `SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE(
    from_account_id = $1
    OR
    to_account_id = $2
) ORDER BY id
LIMIT $3
OFFSET $4
`
	rows, err := q.db.QueryContext(ctx, query,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
