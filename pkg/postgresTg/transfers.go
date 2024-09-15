package postgrestg

import (
	"context"
	"database/sql"
	"tgBank/models"
)

type transferTgSQL struct {
	db *sql.DB
}

func NewTransferSQL(db *sql.DB) TransferTgSQL {
	return &transferTgSQL{
		db: db,
	}
}

func (q *transferTgSQL) CreateTransfer(ctx context.Context, arg models.CreateTransferTgParams) (models.TransferTg, error) {
	query := `INSERT INTO transfers(
    from_chat_id,
    to_chat_id,
    amount
)VALUES(
    $1, $2, $3
) RETURNING id, from_chat_id, to_chat_id, amount, created_at`
	row := q.db.QueryRowContext(ctx, query, arg.FromChatId, arg.ToChatId, arg.Amount)
	var i models.TransferTg
	err := row.Scan(
		&i.ID,
		&i.FromChatId,
		&i.ToChatId,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

func (q *transferTgSQL) GetTransfer(ctx context.Context, id int64) (models.TransferTg, error) {
	query := `SELECT id, from_chat_id, to_chat_id, amount, created_at FROM transfers
WHERE id = $1
LIMIT 1`
	row := q.db.QueryRowContext(ctx, query, id)
	var i models.TransferTg
	err := row.Scan(
		&i.ID,
		&i.FromChatId,
		&i.ToChatId,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

func (q *transferTgSQL) ListTransfers(ctx context.Context, arg models.ListTransfersTgParams) ([]models.TransferTg, error) {
	query := `SELECT id, from_chat_id, to_chat_id, amount, created_at FROM transfers
WHERE(
    from_chat_id = $1
    OR
    to_chat_id = $2
) ORDER BY id
LIMIT $3
OFFSET $4
`
	rows, err := q.db.QueryContext(ctx, query,
		arg.FromChatId,
		arg.ToChatId,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.TransferTg{}
	for rows.Next() {
		var i models.TransferTg
		if err := rows.Scan(
			&i.ID,
			&i.FromChatId,
			&i.ToChatId,
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
