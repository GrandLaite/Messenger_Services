package repository

import (
	"database/sql"
)

type Message struct {
	ID                int
	SenderUsername    string
	RecipientUsername string
	Content           string
	CreatedAt         string
}

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(d *sql.DB) *MessageRepository {
	return &MessageRepository{db: d}
}

func (r *MessageRepository) Create(senderUsername, recipientUsername, content string) (Message, error) {
	q := `
INSERT INTO messages (sender_nickname, recipient_nickname, content)
VALUES ($1, $2, $3)
RETURNING id, sender_nickname, recipient_nickname, content, created_at
`
	var m Message
	err := r.db.QueryRow(q, senderUsername, recipientUsername, content).
		Scan(&m.ID, &m.SenderUsername, &m.RecipientUsername, &m.Content, &m.CreatedAt)
	return m, err
}

func (r *MessageRepository) GetByID(id int) (Message, error) {
	q := `
SELECT id, sender_nickname, recipient_nickname, content, created_at
FROM messages
WHERE id = $1
`
	var m Message
	err := r.db.QueryRow(q, id).Scan(
		&m.ID,
		&m.SenderUsername,
		&m.RecipientUsername,
		&m.Content,
		&m.CreatedAt,
	)
	return m, err
}

// УДАЛЕНО: ListAll() (возврат всех сообщений в системе)

func (r *MessageRepository) Delete(id int) error {
	q := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(q, id)
	return err
}

func (r *MessageRepository) InsertLike(msgID int, userUsername string) error {
	q := `INSERT INTO likes (message_id, user_nickname) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(q, msgID, userUsername)
	return err
}

func (r *MessageRepository) RemoveLike(msgID int, userUsername string) error {
	q := `DELETE FROM likes WHERE message_id = $1 AND user_nickname = $2`
	_, err := r.db.Exec(q, msgID, userUsername)
	return err
}

func (r *MessageRepository) InsertSuperlike(msgID int, userUsername string) error {
	q := `INSERT INTO superlikes (message_id, user_nickname) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(q, msgID, userUsername)
	return err
}

func (r *MessageRepository) RemoveSuperlike(msgID int, userUsername string) error {
	q := `DELETE FROM superlikes WHERE message_id = $1 AND user_nickname = $2`
	_, err := r.db.Exec(q, msgID, userUsername)
	return err
}

func (r *MessageRepository) GetConversation(user1, user2 string) ([]Message, error) {
	q := `
SELECT id, sender_nickname, recipient_nickname, content, created_at
FROM messages
WHERE (sender_nickname = $1 AND recipient_nickname = $2)
   OR (sender_nickname = $2 AND recipient_nickname = $1)
ORDER BY created_at ASC
`
	rows, err := r.db.Query(q, user1, user2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		err := rows.Scan(&m.ID, &m.SenderUsername, &m.RecipientUsername, &m.Content, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *MessageRepository) GetDialogs(username string) ([]string, error) {
	q := `
SELECT DISTINCT partner FROM (
    SELECT recipient_nickname AS partner FROM messages WHERE sender_nickname = $1
    UNION
    SELECT sender_nickname AS partner FROM messages WHERE recipient_nickname = $1
) AS sub
`
	rows, err := r.db.Query(q, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var partners []string
	for rows.Next() {
		var partner string
		err := rows.Scan(&partner)
		if err != nil {
			return nil, err
		}
		partners = append(partners, partner)
	}
	return partners, nil
}
