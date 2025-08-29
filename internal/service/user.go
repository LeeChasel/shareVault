package service

import (
	"github.com/LeeChasel/shareVault/internal/models"
	repoInterface "github.com/LeeChasel/shareVault/internal/repository/interfaces"
	serviceInterface "github.com/LeeChasel/shareVault/internal/service/interfaces"
	"github.com/google/uuid"
)

type userService struct {
	userRepo repoInterface.UserRepository
}

func NewUserService(u repoInterface.UserRepository) serviceInterface.UserService {
	return &userService{
		userRepo: u,
	}
}

func (s *userService) Create(user *models.User) error {
	return s.userRepo.Create(user)
}

func (s *userService) ExistsByEmail(email string) (bool, error) {
	return s.userRepo.ExistsByEmail(email)
}

func (s *userService) ExistsByUsername(username string) (bool, error) {
	return s.userRepo.ExistsByUsername(username)
}

func (s *userService) ExistsByUserId(userId uuid.UUID) (bool, error) {
	return s.userRepo.ExistsByUserId(userId)
}

func (s *userService) FindByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}
