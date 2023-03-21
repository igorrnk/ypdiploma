package model

import "context"

type Servicer interface {
	Registration(*User) (*Token, error)
	Login(*User) (*Token, error)
	AddOrder(*Token, string) error
	GetOrders(*Token) ([]*Order, error)
	GetBalance(*Token) (*Balance, error)
	AddWithdraw(*Token, string, float64) error
	GetWithdraw(*Token) ([]*Withdraw, error)
}

type Repository interface {
	IsUser(context.Context, *User) (bool, error)
	AddUser(context.Context, *User) error
	GetUser(context.Context, *User) error

	//AddToken(token string) error
	//IsToken(token string) (login string, err error)

	//AddOrder() error
	//GetOrders() error
	//UpdateOrder() error
}

type Checker interface {
	Check(string) bool
}
