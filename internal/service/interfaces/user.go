package interfaces

import (
	"github.com/LeeChasel/shareVault/internal/models"
)

type UserService interface {
	Create(user *models.User) error
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
	ExistsByUserId(userId string) (bool, error)
	FindByEmail(email string) (*models.User, error)
}
