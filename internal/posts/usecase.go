package posts

import "main/internal/models"

type PostsUse interface {
	GetPostData(int,[]string) (models.AllPostData, error)
	UpdatePost(int64,string) (models.Message, error)
}