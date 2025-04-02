package service

import (
	"errors"
	"message-service/internal/repository"
)

type MessageService struct {
	msgRepo *repository.MessageRepository
}

func NewMessageService(m *repository.MessageRepository) *MessageService {
	return &MessageService{msgRepo: m}
}

func (s *MessageService) Create(senderUsername, recipientUsername, content string) (repository.Message, error) {
	return s.msgRepo.Create(senderUsername, recipientUsername, content)
}

// УДАЛЕНО: ListAll()

func (s *MessageService) GetByID(id int) (repository.Message, error) {
	return s.msgRepo.GetByID(id)
}

func (s *MessageService) Delete(id int, userUsername string) error {
	message, err := s.msgRepo.GetByID(id)
	if err != nil {
		return err
	}
	if message.SenderUsername != userUsername {
		return errors.New("forbidden")
	}
	return s.msgRepo.Delete(id)
}

func (s *MessageService) LikeMessage(msgID int, userUsername string) error {
	return s.msgRepo.InsertLike(msgID, userUsername)
}

func (s *MessageService) UnlikeMessage(msgID int, userUsername string) error {
	return s.msgRepo.RemoveLike(msgID, userUsername)
}

func (s *MessageService) SuperlikeMessage(msgID int, userUsername string) error {
	return s.msgRepo.InsertSuperlike(msgID, userUsername)
}

func (s *MessageService) UnsuperlikeMessage(msgID int, userUsername string) error {
	return s.msgRepo.RemoveSuperlike(msgID, userUsername)
}

func (s *MessageService) GetConversation(user1, user2 string) ([]repository.Message, error) {
	return s.msgRepo.GetConversation(user1, user2)
}

func (s *MessageService) GetDialogs(username string) ([]string, error) {
	return s.msgRepo.GetDialogs(username)
}
