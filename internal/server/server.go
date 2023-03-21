package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/igorrnk/ypdiploma.git/internal/configs"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
)

type Server struct {
	http.Server
	config   *configs.ServerConfigType
	servicer model.Servicer
	router   chi.Router
	ctx      context.Context
}

func NewServer(config *configs.ServerConfigType, servicer model.Servicer) *Server {
	server := &Server{
		config:   config,
		servicer: servicer,
		router:   chi.NewRouter(),
	}
	server.ctx = context.Background()
	server.Addr = config.ServerAddress
	server.Handler = server.router
	server.router.Post("/api/user/register", server.PostRegister)
	server.router.Post("/api/user/login", server.PostLogin)
	server.router.Post("/api/user/orders", server.PostOrders)
	server.router.Get("/api/user/orders", server.GetOrders)
	return server
}

func (server *Server) Run() error {
	log.Info().Msg("HTTP Server is running.")
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := server.Shutdown(server.ctx); err != nil {
			log.Printf("Service.Run: http.Service.Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	defer server.ctx.Done()
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	<-idleConnsClosed
	server.Close()
	log.Info().Msg("HTTP Server has been stopped.")
	return nil
}
