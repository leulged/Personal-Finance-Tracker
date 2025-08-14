package repositories

import "personal-finance-tracker/domain/entities"

type UserRepository interface {
	CreateUser(user *entities.User) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByID(id string) (*entities.User, error)
	UpdateUser(id string, user *entities.User) (*entities.User, error)
	DeleteUser(id string) error
	CountUsers() (int64, error)
}
