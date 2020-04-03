package usecase

import(
	"main/internal/models"
	"main/internal/forums"
	"main/internal/forums/repository"
)

type ForumUseRealisation struct{
	forumRepo forums.ForumRepo
}

func NewForumUseCaseRealisation(fFrepo repository.ForumRepoRealisation) ForumUseRealisation {
	return ForumUseRealisation{forumRepo:fFrepo}
}

func (ForumU ForumUseRealisation) CreateForum(forum models.Forum) (models.Forum, error){
	return ForumU.forumRepo.CreateNewForum(forum)
}

func (ForumU ForumUseRealisation) GetForumData(slug string) (models.Forum, error){
	return ForumU.forumRepo.GetForum(slug)
}

func (ForumU ForumUseRealisation) CreateThread(slug string, thread models.Thread) (models.Thread, error){
	thread.Forum = slug
	return ForumU.forumRepo.CreateThread(thread)
}

func (ForumU ForumUseRealisation) GetThreads(slug string, limit int , since string , sort bool) ([]models.Thread, error){
	forum := models.Forum{Slug:slug}

	return ForumU.forumRepo.GetThreads(forum,limit,since,sort)
}

func (ForumU ForumUseRealisation) GetForumUsers(slug string, limit int, since string, desc bool) ([]models.UserModel ,error) {
	return ForumU.forumRepo.GetForumUsers(slug,limit,since,desc)
}
