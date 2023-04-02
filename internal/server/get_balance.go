package server

import (
	"encoding/json"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"net/http"
)

func (s *Server) GetBalance(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User").(*model.User)
	balance, err := s.servicer.GetBalance(r.Context(), user)
	var body []byte
	if body, err = json.Marshal(balance); err != nil {
		http.Error(w, "Unable to decode balance.", http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, "Unable to get balance.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)

}
