package server

import (
	"encoding/json"
	"errors"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"net/http"
)

func (s *Server) GetWithDrawals(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User").(*model.User)
	withdraws, err := s.servicer.GetWithdraws(r.Context(), user)
	if errors.Is(err, model.ErrNoWithdraws) {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("No withdraws"))
		return
	}
	var body []byte
	if body, err = json.Marshal(withdraws); err != nil {
		http.Error(w, "Unable to decode withdraws.", http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, "Unable to get withdraws.", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
}
