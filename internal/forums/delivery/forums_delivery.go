package delivery

import (
	"database/sql"
	"github.com/labstack/echo"
	"main/internal/forums"
	"main/internal/forums/usecase"
	"main/internal/models"
	"net/http"
	"strconv"
)

type ForumDelivery struct{
	forumLogic forums.ForumUse
}

func NewForumDeliveryRealisation(fLogic usecase.ForumUseRealisation) ForumDelivery {
	return ForumDelivery{forumLogic:fLogic}
}

func (ForumD ForumDelivery) CreateForum(rwContext echo.Context) error {
	newForumData := new(models.Forum)

	rwContext.Bind(newForumData)

	answer , err := ForumD.forumLogic.CreateForum(*newForumData)

	if err == sql.ErrNoRows {
		return rwContext.JSON(http.StatusNotFound, &models.Error{Message:"Can't find user by nickname: " +newForumData.User})
	}

	if err != nil {
		return rwContext.JSON(http.StatusConflict, answer)
	}

	return rwContext.JSON(http.StatusCreated, answer)

}

func (ForumD ForumDelivery) GetForumUsers(rwContext echo.Context) error {
	slug := rwContext.Param("slug")
	limit , _ := strconv.Atoi(rwContext.QueryParam("limit"))
	since := rwContext.QueryParam("since")
	desc , _:= strconv.ParseBool(rwContext.QueryParam("desc"))

	data , err := ForumD.forumLogic.GetForumUsers(slug,limit,since,desc)

	if err != nil {
		return rwContext.JSON(http.StatusNotFound, &models.Error{Message:"can't find forum by slug: "+slug})
	}

	return rwContext.JSON(http.StatusOK,data)
}

func (ForumD ForumDelivery) GetForum(rwContext echo.Context) error {
	slug := rwContext.Param("slug")

	answer , err := ForumD.forumLogic.GetForumData(slug)

	if err == sql.ErrNoRows{
		return rwContext.JSON(http.StatusNotFound, &models.Error{Message:"Can't find forum by slug: " +slug})
	}

	return rwContext.JSON(http.StatusOK, answer)
}

func (ForumD ForumDelivery) CreateThread(rwContext echo.Context) error {
	slug := rwContext.Param("slug")

	threadReq  := new(models.Thread)

	rwContext.Bind(threadReq)

	thread , err := ForumD.forumLogic.CreateThread(slug, *threadReq)

	if err == sql.ErrNoRows {
		return rwContext.JSON(http.StatusNotFound, &models.Error{Message:"Can't create thread by slug: " + slug })
	}

	if err != nil {
		return rwContext.JSON(http.StatusConflict, thread)
	}

	return rwContext.JSON(http.StatusCreated,thread)
}

func (ForumD ForumDelivery) GetSortedThreads(rwContext echo.Context) error {
	slug := rwContext.Param("slug")
	limit , _:= strconv.Atoi(rwContext.QueryParam("limit"))
	since := rwContext.QueryParam("since")
	desc , _:= strconv.ParseBool(rwContext.QueryParam("desc"))

	threads , err := ForumD.forumLogic.GetThreads(slug,limit,since,desc)

	if err == sql.ErrNoRows {
		return rwContext.JSON(http.StatusNotFound, &models.Error{Message:"Can't create thread by slug: " + slug })
	}

	return rwContext.JSON(http.StatusOK,threads)
}

func (ForumD ForumDelivery) SetupHandlers(server *echo.Echo){
	server.POST("/forum/create", ForumD.CreateForum)
	server.GET("/forum/:slug/details", ForumD.GetForum)
	server.POST("/forum/:slug/create", ForumD.CreateThread)
	server.GET("/forum/:slug/threads", ForumD.GetSortedThreads)
	server.GET("/forum/:slug/users", ForumD.GetForumUsers)
}