package repositories

import (
	"github.com/justadoll/CHAOS/entities"
	"time"
)

type Device interface {
	Insert(device entities.Device) error
	Update(device entities.Device) error
	FindByMacAddress(address string) (*entities.Device, error)
	FindAll(updatedAt time.Time) ([]entities.Device, error)
}
