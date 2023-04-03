package model

import (
	"encoding/json"
	"time"
)

type Withdraw struct {
	Order       string
	Sum         int32
	ProcessedAt time.Time
}

type withdrawJSON struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func (withdraw *Withdraw) UnmarshalJSON(bytes []byte) error {
	alias := &withdrawJSON{}
	if err := json.Unmarshal(bytes, alias); err != nil {
		return err
	}
	withdraw.Order = alias.Order
	withdraw.Sum = int32(alias.Sum * 100)
	return nil
}

func (withdraw *Withdraw) MarshalJSON() ([]byte, error) {
	alias := &withdrawJSON{
		Order:       withdraw.Order,
		Sum:         float64(withdraw.Sum) / 100,
		ProcessedAt: withdraw.ProcessedAt.Format(time.RFC3339),
	}
	return json.Marshal(alias)
}
