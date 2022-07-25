package repositories

import "github.com/justadoll/CHAOS/entities"

type Auth interface {
	Insert(auth entities.Auth) error
	Update(auth entities.Auth) error
	GetFirst() (entities.Auth, error)
}
