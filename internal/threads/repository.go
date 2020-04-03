package threads

import (
	"main/internal/models"
)

type ThreadsRepo interface {
	CreatePost(string, int, []models.Message) ([]models.Message, error)
	VoteThread(string, int, int, models.Thread) (models.Thread, error)
	GetThread(int, models.Thread) (models.Thread, error)
	GetPostsSorted(string, int, int, int, string, bool) ([]models.Message, error)
	UpdateThread(string, int, models.Thread) (models.Thread, error)
}
