package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/justadoll/CHAOS/client/app/gateway"
	"github.com/justadoll/CHAOS/client/app/shared/environment"
)

type ClientGateway struct {
	Configuration *environment.Configuration
	HttpClient    *http.Client
}

func NewGateway(configuration *environment.Configuration, httpClient *http.Client) gateway.Gateway {
	return &ClientGateway{
		Configuration: configuration,
		HttpClient:    httpClient,
	}
}

func (c ClientGateway) NewRequest(method string, url string, body []byte) (*gateway.HttpResponse, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set(c.Configuration.Connection.ContentTypeHeader, c.Configuration.Connection.ContentTypeJSON)
	req.Header.Set(c.Configuration.Connection.CookieHeader, c.Configuration.Connection.Token)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("failed with status code %d", res.StatusCode)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &gateway.HttpResponse{
		ResponseBody: bodyBytes,
		StatusCode:   res.StatusCode,
	}, nil
}
