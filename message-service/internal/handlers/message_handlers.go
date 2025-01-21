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
	sid := r.FormValue("sender_id")
	rid := r.FormValue("recipient_id")
	c := r.FormValue("content")
	senderID, err := strconv.Atoi(sid)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	recipientID, err := strconv.Atoi(rid)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	msg, err := h.srv.Create(senderID, recipientID, c)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
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
	ms, err := h.srv.ListAll()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ms)
}

func (h *MessageHandlers) GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	m, err := h.srv.GetByID(id)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func (h *MessageHandlers) UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	var req struct {
		UserID  int    `json:"user_id"`
		Content string `json:"content"`
		OwnerID int    `json:"owner_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	err = h.srv.Update(id, req.UserID, req.Content, req.OwnerID)
	if err != nil {
		http.Error(w, "", http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	var req struct {
		UserID  int `json:"user_id"`
		OwnerID int `json:"owner_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	err = h.srv.Delete(id, req.UserID, req.OwnerID)
	if err != nil {
		http.Error(w, "", http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) LikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	var req struct {
		UserID int `json:"user_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	err = h.srv.LikeMessage(id, req.UserID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MessageHandlers) SuperlikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	var req struct {
		UserID int `json:"user_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	err = h.srv.SuperlikeMessage(id, req.UserID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
