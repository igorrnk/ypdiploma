package model

import "time"

type StatusType int8

const (
	NewStatus StatusType = iota
	ProcessingStatus
	InvalidStatus
	Processed
)

type Order struct {
	Number     string     `json:"number"`
	Status     StatusType `json:"status"`
	Accrual    float64    `json:"accrual"`
	UploadedAt time.Time  `json:"uploaded_at"`
}

type Balance struct {
	Current  float64 `json:"current"`
	Withdraw float64 `json:"withdraw"`
}

type Withdraw struct {
	Order string
}

type User struct {
	ID       string
	Login    string `json:"login"`
	Password string `json:"password"`
	Hash     string
}

type Token struct {
	Token string
}

func NewToken(user *User) *Token {
	return &Token{Token: user.ID}
}
