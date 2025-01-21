package service

import (
	"errors"
	"message-service/internal/repository"
)

type MessageService struct {
	msgRepo *repository.MessageRepository
	attRepo *repository.AttachmentRepository
}

func NewMessageService(m *repository.MessageRepository, a *repository.AttachmentRepository) *MessageService {
	return &MessageService{msgRepo: m, attRepo: a}
}

func (s *MessageService) Create(senderID, recipientID int, content string) (repository.Message, error) {
	return s.msgRepo.Create(senderID, recipientID, content)
}

func (s *MessageService) ListAll() ([]repository.Message, error) {
	return s.msgRepo.ListAll()
}

func (s *MessageService) GetByID(id int) (repository.Message, error) {
	return s.msgRepo.GetByID(id)
}

func (s *MessageService) Update(id int, userID int, content string, originalOwner int) error {
	if userID != originalOwner {
		return errors.New("forbidden")
	}
	return s.msgRepo.Update(id, content)
}

func (s *MessageService) Delete(id int, userID int, originalOwner int) error {
	if userID != originalOwner {
		return errors.New("forbidden")
	}
	return s.msgRepo.Delete(id)
}

func (s *MessageService) LikeMessage(msgID, userID int) error {
	return s.msgRepo.InsertLike(msgID, userID)
}

func (s *MessageService) SuperlikeMessage(msgID, userID int) error {
	return s.msgRepo.InsertSuperlike(msgID, userID)
}

func (s *MessageService) CreateAttachment(msgID int, data []byte, fType string, fSize int) (repository.Attachment, error) {
	return s.attRepo.Create(msgID, data, fType, fSize)
}

func (s *MessageService) GetAttachments(msgID int) ([]repository.Attachment, error) {
	return s.attRepo.GetByMessageID(msgID)
}
