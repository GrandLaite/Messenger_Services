package repository

import (
	"database/sql"
)

type Attachment struct {
	ID        int
	MessageID int
	FileData  []byte
	FileType  string
	FileSize  int
}

type AttachmentRepository struct {
	db *sql.DB
}

func NewAttachmentRepository(d *sql.DB) *AttachmentRepository {
	return &AttachmentRepository{db: d}
}

func (r *AttachmentRepository) Create(msgID int, data []byte, fType string, fSize int) (Attachment, error) {
	q := `INSERT INTO attachments (message_id, file_data, file_type, file_size) 
	      VALUES ($1, $2, $3, $4) RETURNING id, message_id, file_data, file_type, file_size`
	var a Attachment
	err := r.db.QueryRow(q, msgID, data, fType, fSize).Scan(
		&a.ID, &a.MessageID, &a.FileData, &a.FileType, &a.FileSize,
	)
	return a, err
}

func (r *AttachmentRepository) GetByMessageID(msgID int) ([]Attachment, error) {
	q := `SELECT id, message_id, file_data, file_type, file_size FROM attachments WHERE message_id = $1`
	rows, err := r.db.Query(q, msgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Attachment
	for rows.Next() {
		var a Attachment
		err = rows.Scan(&a.ID, &a.MessageID, &a.FileData, &a.FileType, &a.FileSize)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, nil
}
