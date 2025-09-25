package service

import "github.com/LeeChasel/shareVault/internal/service/interfaces"

type ApplicationServices struct {
	UserService interfaces.UserService
	FileService interfaces.FileService
}
