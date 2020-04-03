package usecase

import (
	"main/internal/models"
	"main/internal/threads"
	"main/internal/threads/repository"
	"strconv"
)

type ThreadsUseRealistaion struct {
	threadRepo threads.ThreadsRepo
}

func NewThreadsUseRealisation(tRep repository.ThreadRepoRealisation) ThreadsUseRealistaion {
	return ThreadsUseRealistaion{threadRepo: tRep}
}

func (ThreadU ThreadsUseRealistaion) CreatePosts(slugOrId string ,posts []models.Message) ([]models.Message, error) {

	id , err := strconv.Atoi(slugOrId)

	if err != nil {
		id = 0
	} else {
		slugOrId = ""
	}

	return ThreadU.threadRepo.CreatePost(slugOrId,id,posts)
}

func (ThreadU ThreadsUseRealistaion) VoteThread(slug , nickname string , voice int) (models.Thread, error){

	threadId , err := strconv.Atoi(slug)

	if err != nil {
		threadId = 0
	} else {
		slug = ""
	}

	return ThreadU.threadRepo.VoteThread(nickname,voice,threadId,models.Thread{Slug:slug})
}

func (ThreadU ThreadsUseRealistaion) GetThread(slug string) (models.Thread, error){

	threadId , err := strconv.Atoi(slug)

	if err != nil {
		threadId = 0
	} else {
		slug = ""
	}

	return ThreadU.threadRepo.GetThread(threadId,models.Thread{Slug:slug})
}

func (ThreadU ThreadsUseRealistaion) GetPosts(slugOrId string,limit int, since int , sortType string , desc bool) ([]models.Message,error) {

	threadId , err := strconv.Atoi(slugOrId)

	if err != nil {
		threadId = 0
	} else {
		slugOrId = ""
	}

	data , err := ThreadU.threadRepo.GetPostsSorted(slugOrId, threadId, limit, since,sortType, desc)
	return data , err
}

func (ThreadU ThreadsUseRealistaion) UpdateThread(slugOrId string, newThreadData models.Thread) (models.Thread, error){
	threadId , err := strconv.Atoi(slugOrId)

	if err != nil {
		threadId = 0
	} else {
		slugOrId = ""
	}

	return ThreadU.threadRepo.UpdateThread(slugOrId, threadId,newThreadData)
}

