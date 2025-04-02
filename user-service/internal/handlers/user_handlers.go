package handlers

import (
	"encoding/json"
	"net/http"
	"user-service/internal/service"

	"github.com/gorilla/mux"
)

type UserHandlers struct {
	srv *service.UserService
}

func NewUserHandlers(s *service.UserService) *UserHandlers {
	return &UserHandlers{srv: s}
}

func (h *UserHandlers) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	id, err := h.srv.CreateUser(req.Username, req.Password, req.Role, req.Email, req.Nickname)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *UserHandlers) CheckPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	u, err := h.srv.CheckPassword(req.Username, req.Password)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"username": u.Username,
		"role":     u.Role,
		"email":    u.Email,
		"nickname": u.Nickname,
	})
}

func (h *UserHandlers) SearchUserHandler(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]
	_, err := h.srv.GetUserByNickname(nickname)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"exists": true})
}
