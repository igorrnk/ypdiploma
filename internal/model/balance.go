package model

import "encoding/json"

type Balance struct {
	Current  int32
	Withdraw int32
}

type balanceJSON struct {
	Current  float64 `json:"current"`
	Withdraw float64 `json:"withdraw"`
}

func (balance *Balance) MarshalJSON() ([]byte, error) {
	alias := &balanceJSON{
		float64(balance.Current) / 100,
		float64(balance.Withdraw) / 100,
	}
	return json.Marshal(alias)
}
