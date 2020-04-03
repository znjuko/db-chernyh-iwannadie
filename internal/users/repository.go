package users

import(
	"main/internal/models"
)

type Repository interface {
	CreateNewUser(models.UserModel) ([]models.UserModel,error)
	UpdateUserData(models.UserModel) (models.UserModel,error)
	GetUserData(string) (models.UserModel,error)
	Status() models.Status
	Clear()
}

