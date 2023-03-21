package model

import "errors"

var (
	ErrDB             = errors.New("database error")
	ErrNoUser         = errors.New("user not found")
	ErrLoginOccupied  = errors.New("login is already occupied")
	ErrWrongLoginPass = errors.New("wrong login or password")
)
