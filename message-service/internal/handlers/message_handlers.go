package handlers

import (
	"encoding/json"
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
	actualSender := r.Header.Get("X-User-Name")
	var req struct {
		RecipientNickname string `json:"recipient_nickname"`
		Content           string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	msg, err := h.srv.Create(actualSender, req.RecipientNickname, req.Content)
	if err != nil {
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// УДАЛЕНО: ListMessagesHandler (получение всех сообщений не требуется)

func (h *MessageHandlers) GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	message, err := h.srv.GetByID(id)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (h *MessageHandlers) DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	actualSender := r.Header.Get("X-User-Name")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	if err := h.srv.Delete(id, actualSender); err != nil {
		http.Error(w, "Forbidden or failed to delete message", http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) LikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	actualUser := r.Header.Get("X-User-Name")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	if err := h.srv.LikeMessage(id, actualUser); err != nil {
		http.Error(w, "Failed to like message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) UnlikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	actualUser := r.Header.Get("X-User-Name")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	if err := h.srv.UnlikeMessage(id, actualUser); err != nil {
		http.Error(w, "Failed to unlike message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) SuperlikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	actualUser := r.Header.Get("X-User-Name")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	if err := h.srv.SuperlikeMessage(id, actualUser); err != nil {
		http.Error(w, "Failed to superlike message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) UnsuperlikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	actualUser := r.Header.Get("X-User-Name")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	if err := h.srv.UnsuperlikeMessage(id, actualUser); err != nil {
		http.Error(w, "Failed to unsuperlike message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) ConversationHandler(w http.ResponseWriter, r *http.Request) {
	actualUser := r.Header.Get("X-User-Name")
	partner := mux.Vars(r)["partner"]
	msgs, err := h.srv.GetConversation(actualUser, partner)
	if err != nil {
		http.Error(w, "Failed to retrieve conversation", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}

func (h *MessageHandlers) DialogsHandler(w http.ResponseWriter, r *http.Request) {
	actualUser := r.Header.Get("X-User-Name")
	dialogs, err := h.srv.GetDialogs(actualUser)
	if err != nil {
		http.Error(w, "Failed to retrieve dialogs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dialogs)
}
