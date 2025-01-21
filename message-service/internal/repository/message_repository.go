package repository

import (
	"database/sql"
)

type Message struct {
	ID          int
	SenderID    int
	RecipientID int
	Content     string
	CreatedAt   string
	UpdatedAt   string
}

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(d *sql.DB) *MessageRepository {
	return &MessageRepository{db: d}
}

func (r *MessageRepository) Create(senderID, recipientID int, content string) (Message, error) {
	q := `INSERT INTO messages (sender_id, recipient_id, content) 
	      VALUES ($1, $2, $3) RETURNING id, sender_id, recipient_id, content, created_at, COALESCE(updated_at, '')`
	var m Message
	err := r.db.QueryRow(q, senderID, recipientID, content).Scan(
		&m.ID, &m.SenderID, &m.RecipientID, &m.Content, &m.CreatedAt, &m.UpdatedAt,
	)
	return m, err
}

func (r *MessageRepository) GetByID(id int) (Message, error) {
	q := `SELECT id, sender_id, recipient_id, content, created_at, COALESCE(updated_at, '') FROM messages WHERE id = $1`
	var m Message
	err := r.db.QueryRow(q, id).Scan(
		&m.ID, &m.SenderID, &m.RecipientID, &m.Content, &m.CreatedAt, &m.UpdatedAt,
	)
	return m, err
}

func (r *MessageRepository) ListAll() ([]Message, error) {
	q := `SELECT id, sender_id, recipient_id, content, created_at, COALESCE(updated_at, '') FROM messages`
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Message
	for rows.Next() {
		var m Message
		err = rows.Scan(
			&m.ID, &m.SenderID, &m.RecipientID, &m.Content, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}

func (r *MessageRepository) Update(id int, content string) error {
	q := `UPDATE messages SET content = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(q, content, id)
	return err
}

func (r *MessageRepository) Delete(id int) error {
	q := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(q, id)
	return err
}

func (r *MessageRepository) InsertLike(msgID, userID int) error {
	q := `INSERT INTO likes (message_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(q, msgID, userID)
	return err
}

func (r *MessageRepository) InsertSuperlike(msgID, userID int) error {
	q := `INSERT INTO superlikes (message_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(q, msgID, userID)
	return err
}
