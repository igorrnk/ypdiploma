package service

import (
	"context"
	"errors"
	"github.com/igorrnk/ypdiploma.git/internal/configs"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"github.com/igorrnk/ypdiploma.git/internal/server"
	"github.com/igorrnk/ypdiploma.git/internal/storage"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
)

// Gophermart is an implementation of the model.servicer interface
type Gophermart struct {
	context    context.Context
	storage    model.Repository
	checker    model.Checker
	httpServer *server.Server
	workerPool *WorkerPool
}

func NewGophermart(ctx context.Context, config *configs.ConfigType) (*Gophermart, error) {
	gophermart := &Gophermart{
		context: ctx,
		//Storage: storage.NewPostgresStorage(ctx, &configs.DBConfigType),
		//server:  server.NewServer(&configs.ServerConfigType),
	}
	gophermart.httpServer = server.NewServer(&config.ServerConfigType, gophermart)
	var err error
	gophermart.storage, err = storage.NewPostgresStorage(ctx, &config.DBConfigType)
	return gophermart, err
}

func (gm *Gophermart) Run() error {

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := gm.httpServer.Shutdown(gm.context); err != nil {
			log.Error().Msgf("Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	defer gm.context.Done()
	log.Info().Msg("HTTP Server is running.")
	if err := gm.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	<-idleConnsClosed
	log.Info().Msg("HTTP Server has been stopped.")
	gm.Close()
	return nil
}

func (gm *Gophermart) Close() {

}

func (gm *Gophermart) Registration(user *model.User) (*model.Token, error) {
	isOccupied, err := gm.storage.IsUser(gm.context, user)
	if err != nil {
		log.Error().Err(err).Msg("database error")
		return nil, err
	}
	if isOccupied {
		return nil, model.ErrLoginOccupied
	}
	user.Hash = GenHash(user.Password)
	err = gm.storage.AddUser(gm.context, user)

	if err != nil {
		return nil, err
	}

	return model.NewToken(user), nil
}

func (gm *Gophermart) Login(user *model.User) (*model.Token, error) {
	err := gm.storage.GetUser(gm.context, user)
	if errors.Is(err, model.ErrNoUser) {
		return nil, model.ErrWrongLoginPass
	}
	if err != nil {
		return nil, model.ErrDB
	}
	return model.NewToken(user), nil
}

func (gm *Gophermart) AddOrder(*model.Token, string) error {
	//TODO implement me
	panic("implement me")
}

func (gm *Gophermart) GetOrders(*model.Token) ([]*model.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (gm *Gophermart) GetBalance(*model.Token) (*model.Balance, error) {
	//TODO implement me
	panic("implement me")
}

func (gm *Gophermart) AddWithdraw(*model.Token, string, float64) error {
	//TODO implement me
	panic("implement me")
}

func (gm *Gophermart) GetWithdraw(*model.Token) ([]*model.Withdraw, error) {
	//TODO implement me
	panic("implement me")
}
