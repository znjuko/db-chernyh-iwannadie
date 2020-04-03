package users

import(
	"main/internal/models"
)

type UserUseCase interface {
	GetUser(string) (models.UserModel , error)
	CreateUser(models.UserModel) (interface{}, error)
	UpdateUserData(models.UserModel) (models.UserModel, error)
	GetServerStatus() models.Status
	Clear()
}
