package interfaces

import (
	"github.com/LeeChasel/shareVault/internal/models"
	"github.com/google/uuid"
)

type UserService interface {
	Create(user *models.User) error
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
	ExistsByUserId(userId uuid.UUID) (bool, error)
	FindByEmail(email string) (*models.User, error)
}
