package server

import (
	"encoding/json"
	"errors"
	"github.com/igorrnk/ypdiploma.git/internal/auth"
	"github.com/igorrnk/ypdiploma.git/internal/model"

	"io"
	"net/http"
)

func (s *Server) PostRegister(w http.ResponseWriter, r *http.Request) {

	user := model.User{}
	var body []byte
	var err error
	if body, err = io.ReadAll(r.Body); err != nil {
		http.Error(w, "Unable to read body.", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &user); err != nil {
		http.Error(w, "Unable to decode body.", http.StatusBadRequest)
		return
	}
	err = s.servicer.Registration(r.Context(), &user)
	if errors.Is(err, model.ErrLoginOccupied) {
		http.Error(w, "Unable to register a new user. The login is already occupied.", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Unable to register a new user.", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Authorization", "Bearer "+auth.TokenByUser(&user))
	w.WriteHeader(http.StatusOK)
}
