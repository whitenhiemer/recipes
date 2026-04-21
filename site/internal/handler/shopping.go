package handler

import (
	"encoding/json"
	"net/http"
)

type shoppingRequest struct {
	Slugs []string `json:"slugs"`
}

func (s *Server) handleShoppingListAPI(w http.ResponseWriter, r *http.Request) {
	var req shoppingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	list := s.idx.GenerateShoppingList(req.Slugs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
