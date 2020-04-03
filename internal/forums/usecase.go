package forums

import (
	"main/internal/models"
)

type ForumUse interface {
	CreateForum(models.Forum) (models.Forum, error)
	GetForumData(string) (models.Forum, error)
	CreateThread(string, models.Thread) (models.Thread, error)
	GetThreads(string, int, string, bool) ([]models.Thread, error)
	GetForumUsers(string, int, string, bool) ([]models.UserModel ,error)
}
