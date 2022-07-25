package device

import (
	"github.com/justadoll/CHAOS/entities"
)

type Service interface {
	Insert(entities.Device) error
	FindAll() ([]entities.Device, error)
}
