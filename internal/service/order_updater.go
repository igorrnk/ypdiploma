package service

import (
	"context"
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
			updater.UpdateOrder(order)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (updater *OrderUpdater) GetOrders(ctx context.Context) ([]*model.Order, error) {
	return updater.Storage.GetAllOrders(ctx)
}

func (updater *OrderUpdater) UpdateOrder(order *model.Order) {
	err := updater.Client.UpdateOrder(order)
	if err != nil {
		return
	}
}
