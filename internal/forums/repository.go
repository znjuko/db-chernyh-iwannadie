package forums

import (
	"main/internal/models"
)

type ForumRepo interface {
	CreateNewForum(models.Forum) (models.Forum, error)
	GetForum(string) (models.Forum, error)
	CreateThread(models.Thread) (models.Thread, error)
	GetThreads(models.Forum,int,string,bool) ([]models.Thread,error)
	GetForumUsers(string, int, string, bool) ([]models.UserModel ,error)
}
