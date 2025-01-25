package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"auth-service/internal/service"
)

type AuthHandlers struct {
	srv *service.AuthService
}

func NewAuthHandlers(s *service.AuthService) *AuthHandlers {
	return &AuthHandlers{srv: s}
}

func (h *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var c struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	token, err := h.srv.LoginUser(c.Username, c.Password)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	url := os.Getenv("USER_SERVICE_URL")
	if url == "" {
		url = "http://localhost:8082"
	}
	buf, _ := json.Marshal(reqData)
	u := url + "/users"
	rq, err := http.NewRequest("POST", u, bytes.NewBuffer(buf))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	rq.Header.Set("Content-Type", "application/json")
	cl := &http.Client{}
	resp, err := cl.Do(rq)
	if err != nil {
		http.Error(w, "Failed to connect to user service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, string(body), resp.StatusCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
