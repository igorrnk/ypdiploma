package model

import "errors"

var (
	ErrDB = errors.New("database error")

	ErrNoUser         = errors.New("user not found")
	ErrLoginOccupied  = errors.New("login is already occupied")
	ErrWrongLoginPass = errors.New("wrong login or password")
	ErrWrongAuth      = errors.New("invalid authentication")

	ErrOrderOccupied = errors.New("order is already occupied")
	ErrOrderUpload   = errors.New("order was uploaded by user")
	ErrOrderNumber   = errors.New("wrong order number")
	ErrNoOrders      = errors.New("no orders")

	ErrNoWithdraws = errors.New("no withdraws")

	ErrInsufFunds = errors.New("insufficient funds")

	ErrAccrual         = errors.New(("accrual error"))
	ErrTooManyRequests = errors.New("high rate")
)
