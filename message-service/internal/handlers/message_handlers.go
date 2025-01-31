package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"message-service/internal/service"

	"github.com/gorilla/mux"
)

type MessageHandlers struct {
	srv *service.MessageService
}

func NewMessageHandlers(s *service.MessageService) *MessageHandlers {
	return &MessageHandlers{srv: s}
}

func (h *MessageHandlers) CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	senderID, _ := strconv.Atoi(r.FormValue("sender_id"))
	recipientID, _ := strconv.Atoi(r.FormValue("recipient_id"))
	content := r.FormValue("content")

	msg, err := h.srv.Create(senderID, recipientID, content)
	if err != nil {
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}

	file, header, _ := r.FormFile("attachment")
	if file != nil && header != nil {
		defer file.Close()
		buf, _ := io.ReadAll(file)
		fSize := len(buf)
		fType := header.Header.Get("Content-Type")
		_, _ = h.srv.CreateAttachment(msg.ID, buf, fType, fSize)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func (h *MessageHandlers) ListMessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := h.srv.ListAll()
	if err != nil {
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *MessageHandlers) GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	message, err := h.srv.GetByID(id)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (h *MessageHandlers) UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req struct {
		Content string `json:"content"`
		UserID  int    `json:"user_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	err := h.srv.Update(id, req.UserID, req.Content)
	if err != nil {
		http.Error(w, "Forbidden or failed to update message", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req struct {
		UserID int `json:"user_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	err := h.srv.Delete(id, req.UserID)
	if err != nil {
		http.Error(w, "Forbidden or failed to delete message", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) LikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req struct {
		UserID int `json:"user_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	err := h.srv.LikeMessage(id, req.UserID)
	if err != nil {
		http.Error(w, "Failed to like message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) UnlikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req struct {
		UserID int `json:"user_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	err := h.srv.UnlikeMessage(id, req.UserID)
	if err != nil {
		http.Error(w, "Failed to unlike message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) SuperlikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req struct {
		UserID int `json:"user_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	err := h.srv.SuperlikeMessage(id, req.UserID)
	if err != nil {
		http.Error(w, "Failed to superlike message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) UnsuperlikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req struct {
		UserID int `json:"user_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	err := h.srv.UnsuperlikeMessage(id, req.UserID)
	if err != nil {
		http.Error(w, "Failed to unsuperlike message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
