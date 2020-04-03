package usecase

import (
	"main/internal/models"
	"main/internal/users"
	"main/internal/users/repository"
)

type UserUseCaseRealisation struct {
	repositoryLogic users.Repository
}

func NewUserUseCaseRealisation(db repository.UserRepoRealisation) UserUseCaseRealisation {
	return UserUseCaseRealisation{repositoryLogic: db}
}

func (UserU UserUseCaseRealisation) GetUser(nickname string) (models.UserModel , error) {

	return UserU.repositoryLogic.GetUserData(nickname)

}

func (UserU UserUseCaseRealisation) CreateUser(newUser models.UserModel) (interface{} , error) {

	answerData , err := UserU.repositoryLogic.CreateNewUser(newUser)

	if err != nil {
		return answerData , err
	}

	return answerData[0] , err

}

func (UserU UserUseCaseRealisation) UpdateUserData(newUserData models.UserModel) (models.UserModel, error) {

	return UserU.repositoryLogic.UpdateUserData(newUserData)
}

func (UserU UserUseCaseRealisation) GetServerStatus() models.Status {
	return UserU.repositoryLogic.Status()
}

func (UserU UserUseCaseRealisation) Clear()  {
	UserU.repositoryLogic.Clear()
}
