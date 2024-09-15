package postgrestg

import (
	"database/sql"
	"tgBank/models"
)

type messageTgSQL struct {
	db *sql.DB
}

func NewMessageTgSQL(db *sql.DB) MessageTgSQL {
	return &messageTgSQL{
		db: db,
	}
}

func (r *messageTgSQL) CreateMessage(msg models.Message) error {
	query := `
    INSERT INTO messages (chat_id, message_id, text)
    VALUES ($1, $2, $3)
    `
	_, err := r.db.Exec(query, msg.ChatId, msg.MessageId, msg.Text)
	return err
}
