package client

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/igorrnk/ypdiploma.git/internal/configs"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"net/http"
)

type RestyClient struct {
	client         *resty.Client
	addressAccrual string
}

func NewRestyClient(config *configs.ConfigType) *RestyClient {
	client := &RestyClient{
		client:         resty.New(),
		addressAccrual: config.AccrualSysAddress,
	}
	return client
}

func (client *RestyClient) GetOrder(order *model.Order) error {
	url := fmt.Sprintf("%s/api/orders/%s", client.addressAccrual, order.Number)

	resp, err := client.client.R().Get(url)
	if err != nil {
		return model.ErrAccrual
	}
	if resp != nil {
		if resp.StatusCode() == http.StatusTooManyRequests {
			return model.ErrTooManyRequests
		}
		if resp.StatusCode() == http.StatusNoContent {
			return nil
		}
		if resp.StatusCode() != http.StatusOK {
			return model.ErrAccrual
		}
		err = json.Unmarshal(resp.Body(), order)
		if err != nil {
			return model.ErrAccrual
		}
	}
	return nil
}
