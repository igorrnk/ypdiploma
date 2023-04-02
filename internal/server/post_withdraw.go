package server

import (
	"encoding/json"
	"errors"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"io"
	"net/http"
)

func (s *Server) PostWithdraw(w http.ResponseWriter, r *http.Request) {

	withdraw := model.Withdraw{}
	var body []byte
	var err error
	if body, err = io.ReadAll(r.Body); err != nil {
		http.Error(w, "Unable to read body.", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &withdraw); err != nil {
		http.Error(w, "Unable to decode body.", http.StatusBadRequest)
		return
	}
	user := r.Context().Value("User").(*model.User)

	err = s.servicer.AddWithdraw(r.Context(), user, &withdraw)
	if errors.Is(err, model.ErrInsufFunds) {
		http.Error(w, "Insufficient funds.", http.StatusPaymentRequired)
		return
	}
	if errors.Is(err, model.ErrOrderNumber) {
		http.Error(w, "Wrong order number.", http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		http.Error(w, "Unable to withdraw.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Withdraw is done."))
}
