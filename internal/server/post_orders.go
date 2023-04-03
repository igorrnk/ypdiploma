package server

import (
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func (s *Server) PostOrders(w http.ResponseWriter, r *http.Request) {
	var body []byte
	var err error
	if body, err = io.ReadAll(r.Body); err != nil {
		http.Error(w, "Unable to read body.", http.StatusBadRequest)
		return
	}
	user := r.Context().Value("User").(*model.User)

	err = s.servicer.AddOrder(r.Context(), user, string(body))

	if errors.Is(err, model.ErrOrderUpload) {
		w.WriteHeader(http.StatusOK)
		return
	}
	if errors.Is(err, model.ErrOrderOccupied) {
		http.Error(w, "Order number was upload another user.", http.StatusConflict)
		return
	}
	if errors.Is(err, model.ErrOrderNumber) {
		http.Error(w, "Wrong order number.", http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		http.Error(w, "Unable to add order.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	return
}
