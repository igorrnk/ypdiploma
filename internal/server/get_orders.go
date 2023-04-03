package server

import (
	"encoding/json"
	"errors"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"net/http"
)

func (s *Server) GetOrders(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User").(*model.User)
	orders, err := s.servicer.GetOrders(r.Context(), user)
	var body []byte
	if body, err = json.Marshal(orders); err != nil {
		http.Error(w, "Unable to decode orders to JSON.", http.StatusInternalServerError)
		return
	}
	if errors.Is(err, model.ErrNoOrders) {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("No orders"))
		return
	}
	if err != nil {
		http.Error(w, "Unable to get orders.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
}
