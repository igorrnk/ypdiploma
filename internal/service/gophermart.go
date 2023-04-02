package service

import (
	"context"
	"errors"
	"github.com/igorrnk/ypdiploma.git/internal/client"
	"github.com/igorrnk/ypdiploma.git/internal/configs"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"github.com/igorrnk/ypdiploma.git/internal/server"
	"github.com/igorrnk/ypdiploma.git/internal/storage"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Gophermart is an implementation of the model.servicer interface
type Gophermart struct {
	context    context.Context
	storage    model.Repository
	checker    model.Checker
	httpServer *server.Server
	updater    *OrderUpdater
}

func NewGophermart(ctx context.Context, config *configs.ConfigType) (*Gophermart, error) {
	gophermart := &Gophermart{
		context: ctx,
		//Storage: storage.NewPostgresStorage(ctx, &configs.DBConfigType),
		//server:  server.NewServer(&configs.ServerConfigType),
		checker: &LuhnChecker{},
	}
	gophermart.httpServer = server.NewServer(&config.ServerConfigType, gophermart)
	var err error
	gophermart.storage, err = storage.NewPostgresStorage(ctx, &config.DBConfigType)
	gophermart.updater = &OrderUpdater{
		Client:  client.NewRestyClient(config),
		Storage: gophermart.storage,
	}
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

	go gm.updater.Run(gm.context)

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

func (gm *Gophermart) Registration(ctx context.Context, user *model.User) error {
	user.Hash = GenHash(user.Password)
	err := gm.storage.AddUser(ctx, user)
	if errors.Is(err, model.ErrLoginOccupied) {
		return model.ErrLoginOccupied
	}
	if err != nil {
		return err
	}
	return nil
}

func (gm *Gophermart) Login(ctx context.Context, user *model.User) error {
	err := gm.storage.GetUser(gm.context, user)
	if errors.Is(err, model.ErrNoUser) {
		return model.ErrWrongLoginPass
	}
	if err != nil {
		return model.ErrDB
	}
	return nil
}

func (gm *Gophermart) AddOrder(ctx context.Context, user *model.User, orderNumber string) error {
	if !gm.checker.Check(orderNumber) {
		return model.ErrOrderNumber
	}
	order := &model.Order{
		Number:     orderNumber,
		Status:     model.NewStatus,
		Accrual:    0,
		UploadedAt: time.Now(),
	}
	err := gm.storage.AddOrder(ctx, user, order)

	if err != nil {
		return err
	}
	return nil
}

func (gm *Gophermart) GetOrders(ctx context.Context, user *model.User) (orders []*model.Order, err error) {
	return gm.storage.GetOrders(ctx, user)
}

func (gm *Gophermart) GetBalance(ctx context.Context, user *model.User) (*model.Balance, error) {
	balance := model.Balance{}
	sum, err := gm.storage.GetSumAccrual(ctx, user)
	if err != nil {
		return nil, err
	}
	withdrawn, err := gm.storage.GetSumWithdrawn(ctx, user)
	if err != nil {
		return nil, err
	}
	balance.Current = sum
	balance.Withdraw = withdrawn
	return &balance, nil
}

func (gm *Gophermart) AddWithdraw(ctx context.Context, user *model.User, withdraw *model.Withdraw) error {
	balance, err := gm.GetBalance(ctx, user)
	if err != nil {
		return model.ErrDB
	}
	if withdraw.Sum > balance.Current-balance.Withdraw {
		return model.ErrInsufFunds
	}
	withdraw.ProcessedAt = time.Now()
	return gm.storage.AddWithdraw(ctx, user, withdraw)
}

func (gm *Gophermart) GetWithdraws(ctx context.Context, user *model.User) ([]*model.Withdraw, error) {
	return gm.storage.GetWithdraws(ctx, user)
}
