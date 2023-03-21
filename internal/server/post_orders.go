package server

import (
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
	auth := r.Header.Get("Authorization")

	s.servicer.AddOrder()

}
