package usecase

import (
	"main/internal/models"
	"main/internal/posts"
	"main/internal/posts/repository"
)

type PostUseRealisation struct {
	postRepo posts.PostRepo
}

func NewPostUseRealistaion(pRep repository.PostRepoRealisation) PostUseRealisation {
	return PostUseRealisation{postRepo: pRep}
}

func (PostU PostUseRealisation) GetPostData(id int, flags []string) (models.AllPostData, error) {
	return PostU.postRepo.GetPost(id, flags)
}

func (PostU PostUseRealisation) UpdatePost(id int64, message string) (models.Message, error) {
	return PostU.postRepo.UpdatePost(models.Message{Id:id,Message:message})
}