package model

import (
	"encoding/json"
	"time"
)

const (
	NewStatus        = "NEW"
	ProcessingStatus = "PROCESSING"
	InvalidStatus    = "INVALID"
	ProcessedStatus  = "PROCESSED"
)

type Order struct {
	Number     string
	Status     string
	Accrual    int32
	UploadedAt time.Time
}

type orderJSON struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

func (order *Order) MarshalJSON() ([]byte, error) {
	alias := &orderJSON{
		Number: order.Number,
		Status: order.Status,
	}
	if order.Status == ProcessedStatus {
		accrual := float64(order.Accrual) / 100
		alias.Accrual = accrual
	}
	uploadedAt := order.UploadedAt.Format(time.RFC3339)
	alias.UploadedAt = uploadedAt
	return json.Marshal(alias)
}

func (order *Order) UnmarshalJSON(bytes []byte) error {
	alias := &orderJSON{}
	if err := json.Unmarshal(bytes, alias); err != nil {
		return err
	}
	order.Number = alias.Number
	order.Status = alias.Status
	order.Accrual = int32(alias.Accrual * 100)
	return nil
}
