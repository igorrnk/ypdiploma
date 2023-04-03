package model

import "context"

type Servicer interface {
	Registration(ctx context.Context, user *User) (err error)
	Login(ctx context.Context, user *User) (err error)
	AddOrder(ctx context.Context, user *User, order string) (err error)
	GetOrders(ctx context.Context, user *User) (orders []*Order, err error)
	GetBalance(ctx context.Context, user *User) (balance *Balance, err error)
	AddWithdraw(ctx context.Context, user *User, withdraw *Withdraw) (err error)
	GetWithdraws(ctx context.Context, user *User) (withdraws []*Withdraw, err error)
}

type Repository interface {
	AddUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, user *User) error

	AddOrder(ctx context.Context, user *User, order *Order) error
	GetOrders(ctx context.Context, user *User) (orders []*Order, err error)

	GetAllOrders(ctx context.Context) (orders []*Order, err error)
	UpdateOrder(ctx context.Context, order *Order) error

	GetSumAccrual(ctx context.Context, user *User) (sum int32, err error)
	GetSumWithdrawn(ctx context.Context, user *User) (sum int32, err error)

	AddWithdraw(ctx context.Context, user *User, withdraw *Withdraw) error
	GetWithdraws(ctx context.Context, user *User) (withdraws []*Withdraw, err error)
}

type Checker interface {
	Check(string) bool
}
