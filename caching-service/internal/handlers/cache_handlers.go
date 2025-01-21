package handlers

import (
	"encoding/json"
	"net/http"

	"caching-service/internal/service"
)

type CacheHandlers struct {
	srv *service.CacheService
}

func NewCacheHandlers(s *service.CacheService) *CacheHandlers {
	return &CacheHandlers{srv: s}
}

func (h *CacheHandlers) SetHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	err = h.srv.SetValue(req.Key, req.Value)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *CacheHandlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	k := r.URL.Query().Get("key")
	val, err := h.srv.GetValue(k)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"value": val})
}
