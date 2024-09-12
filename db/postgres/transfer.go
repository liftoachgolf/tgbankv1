package postgres

import (
	"context"
	"database/sql"
	"tgBank/models"
)

type TransferSQLIMPl struct {
	db *sql.DB
}

func NewTransferSQL(db *sql.DB) *TransferSQLIMPl {
	return &TransferSQLIMPl{
		db: db,
	}
}

func (q *TransferSQLIMPl) CreateTransfer(ctx context.Context, arg models.CreateTransferParams) (models.Transfer, error) {
	query := `INSERT INTO transfers(
    from_account_id,
    to_account_id,
    amount
)VALUES(
    $1, $2, $3
) RETURNING id, from_account_id, to_account_id, amount, created_at`
	row := q.db.QueryRowContext(ctx, query, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i models.Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

func (q *TransferSQLIMPl) GetTransfer(ctx context.Context, id int64) (models.Transfer, error) {
	query := `SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE id = $1
LIMIT 1`
	row := q.db.QueryRowContext(ctx, query, id)
	var i models.Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

func (q *TransferSQLIMPl) ListTransfers(ctx context.Context, arg models.ListTransfersParams) ([]models.Transfer, error) {
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
	items := []models.Transfer{}
	for rows.Next() {
		var i models.Transfer
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
