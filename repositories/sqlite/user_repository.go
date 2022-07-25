package sqlite

import (
	"github.com/jinzhu/gorm"
	"github.com/justadoll/CHAOS/entities"
	"github.com/justadoll/CHAOS/repositories"
)

type userSqliteRepository struct {
	dbClient *gorm.DB
}

func NewUserRepository(database *gorm.DB) repositories.User {
	return &userSqliteRepository{dbClient: database}
}

func (u userSqliteRepository) Insert(user entities.User) error {
	return u.dbClient.Create(&user).Error
}

func (u userSqliteRepository) Update(user *entities.User) error {
	return u.dbClient.Model(&user).Where(
		entities.User{Username: user.Username}).Update(&user).Error
}

func (u userSqliteRepository) FindByUsername(username string) (*entities.User, error) {
	var user entities.User
	if err := u.dbClient.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, handleError(err)
	}
	return &user, nil
}
