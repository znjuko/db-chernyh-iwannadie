package threads

import "main/internal/models"

type ThreadsUse interface {
	CreatePosts(string ,[]models.Message) ([]models.Message, error)
	VoteThread(string, string, int) (models.Thread, error)
	GetThread(string) (models.Thread, error)
	GetPosts(string,int,int  , string, bool) ([]models.Message,error)
	UpdateThread(string,models.Thread) (models.Thread, error)
}