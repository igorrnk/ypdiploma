package server

import (
	"encoding/json"
	"errors"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"io"
	"net/http"
)

func (s *Server) PostLogin(w http.ResponseWriter, r *http.Request) {
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
	token, err := s.servicer.Login(&user)
	if errors.Is(err, model.ErrWrongLoginPass) {
		http.Error(w, "Wrong login or password.", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Unable to register a new user.", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Authorization", token.Token)
	w.WriteHeader(http.StatusOK)

}
