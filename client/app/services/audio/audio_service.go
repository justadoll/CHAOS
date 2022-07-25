package audio

import (
	"fmt"

	"github.com/justadoll/CHAOS/client/app/services"
)

type AudioService struct{}

func NewAudioService() services.Audio {
	return &AudioService{}
}

func (d AudioService) Record(raw_seconds string) ([]byte, error) {
	tmp := make([]byte, 4)
	fmt.Println("Recording seconds: ", raw_seconds)
	return tmp, nil
}
