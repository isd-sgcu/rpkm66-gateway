package vaccine

import (
	"github.com/go-resty/resty/v2"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/config"
	"github.com/pkg/errors"
	"net/http"
)

type Client struct {
	client *resty.Client
}

func NewClient(conf config.Vaccine) *Client {
	client := resty.New().
		SetHeader("Authorization", conf.ApiKey).
		SetBaseURL(conf.Host)

	return &Client{client: client}
}

func (c *Client) Verify(req *dto.VaccineRequest, res *dto.VaccineResponse) error {
	resp, err := c.client.R().
		SetBody(req).
		SetResult(res).
		Post("/vaccine")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		return errors.New("Invalid QR")
	}

	return nil
}
