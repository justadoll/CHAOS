package device

import (
	"errors"
	"github.com/justadoll/CHAOS/entities"
	"github.com/justadoll/CHAOS/internal/utils"
	"github.com/justadoll/CHAOS/repositories"
	"time"
)

type deviceService struct {
	Repository repositories.Device
}

func NewDeviceService(repository repositories.Device) Service {
	return &deviceService{Repository: repository}
}

func (d deviceService) Insert(input entities.Device) error {
	_, err := d.Repository.FindByMacAddress(input.MacAddress)
	if errors.Is(err, repositories.ErrNotFound) {
		return d.Repository.Insert(input)
	}
	return d.Repository.Update(input)
}

func (d deviceService) FindAll() ([]entities.Device, error) {
	devices, err := d.Repository.FindAll(time.Now().Add(time.Minute * time.Duration(-3)))
	if err != nil {
		return nil, err
	}
	for index, device := range devices {
		devices[index].MacAddressBase64 = utils.EncodeBase64(device.MacAddress)
	}
	return devices, nil
}
