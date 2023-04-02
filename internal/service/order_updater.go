package service

import (
	"context"
	"errors"
	"github.com/igorrnk/ypdiploma.git/internal/client"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"time"
)

type OrderUpdater struct {
	Client  *client.RestyClient
	Storage model.Repository
}

func (updater *OrderUpdater) Run(ctx context.Context) {
	for {
		orders, _ := updater.GetOrders(ctx)
		for _, order := range orders {
			updater.UpdateOrder(ctx, order)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (updater *OrderUpdater) GetOrders(ctx context.Context) ([]*model.Order, error) {
	return updater.Storage.GetAllOrders(ctx)
}

func (updater *OrderUpdater) UpdateOrder(ctx context.Context, order *model.Order) {
	err := updater.Client.UpdateOrder(order)
	if errors.Is(err, model.ErrTooManyRequests) {
		time.Sleep(time.Second)
		return
	}
	if err != nil {
		return
	}
	err = updater.Storage.UpdateOrder(ctx, order)
	if err != nil {
		return
	}
}
