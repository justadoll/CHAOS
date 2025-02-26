package url

import (
	"context"
	"fmt"
	"github.com/justadoll/CHAOS/services/client"
	"net/url"
	"strings"
)

type urlService struct {
	client.Service
}

func NewUrlService(clientService client.Service) Service {
	return &urlService{clientService}
}

func (u urlService) OpenUrl(ctx context.Context, address string, rawUrl string) error {
	if !strings.Contains(strings.ToLower(rawUrl), "http") {
		rawUrl = fmt.Sprintf("https://%s", rawUrl)
	}
	urlStr, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}
	if _, err := u.SendCommand(ctx, client.SendCommandInput{
		MacAddress: address,
		Request:    fmt.Sprintf("open-url %s", urlStr.String()),
	}); err != nil {
		return err
	}
	return nil
}
