package posts

import(
	"main/internal/models"
)

type PostRepo interface {
	GetPost(int,[]string) (models.AllPostData, error)
	UpdatePost(models.Message) (models.Message, error)
}
