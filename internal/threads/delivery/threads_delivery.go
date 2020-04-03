package delivery

import (
	"database/sql"
	"github.com/labstack/echo"
	"main/internal/models"
	"main/internal/threads"
	"main/internal/threads/usecase"
	"net/http"
	"strconv"
)

type ThreadDelivery struct {
	threadLogic threads.ThreadsUse
}

func NewThreadDeliveryRealisation(fLogic usecase.ThreadsUseRealistaion) ThreadDelivery {
	return ThreadDelivery{threadLogic: fLogic}
}

func (Thread ThreadDelivery) CreatePosts(rwContext echo.Context) error {

	slugOrId := rwContext.Param("slug_or_id")

	posts := []models.Message{}
	rwContext.Bind(&posts)

	posts, err := Thread.threadLogic.CreatePosts(slugOrId, posts)

	if err == sql.ErrNoRows {
		return rwContext.JSON(http.StatusNotFound, models.Error{Message: "can't find thread by slug_or_id: " + slugOrId})
	}

	if err != nil {
		if err.Error() == "no user" {
			return rwContext.JSON(http.StatusNotFound, models.Error{Message:"Can't find post author by nickname: " + posts[0].Author})

		}

		if err.Error() == "No parent message!"{
			return rwContext.JSON(http.StatusConflict, models.Error{Message:err.Error()})
		}

		if err.Error() == "Parent post was created in another thread"{
			return rwContext.JSON(http.StatusConflict, models.Error{Message:err.Error()})
		}
	}


	if err != nil {
		return rwContext.JSON(http.StatusConflict, models.Error{Message: "no parent message"})
	}

	return rwContext.JSON(http.StatusCreated, posts)

}

func (Thread ThreadDelivery) VoteThread(rwContext echo.Context) error {
	slugOrId := rwContext.Param("slug_or_id")

	voice := new(models.Voice)
	rwContext.Bind(&voice)

	thread , err := Thread.threadLogic.VoteThread(slugOrId,voice.Nickname,voice.Voice)

	if err != nil {
		return rwContext.JSON(http.StatusNotFound, models.Error{Message:"can't vote by slug_or_id:" + slugOrId})
	}

	return rwContext.JSON(http.StatusOK,thread)
}

func (Thread ThreadDelivery) GetThread(rwContext echo.Context) error {
	slugOrId := rwContext.Param("slug_or_id")

	thread , err := Thread.threadLogic.GetThread(slugOrId)

	if err != nil {
		return rwContext.JSON(http.StatusNotFound, models.Error{Message:"can't get thread by slug_or_id:" + slugOrId})
	}

	return rwContext.JSON(http.StatusOK,thread)
}

func (Thread ThreadDelivery) GetPosts(rwContext echo.Context) error {
	slugOrId := rwContext.Param("slug_or_id")
	limit , _ := strconv.Atoi(rwContext.QueryParam("limit"))
	since , _:= strconv.Atoi(rwContext.QueryParam("since"))
	sortType := rwContext.QueryParam("sort")
	desc , err := strconv.ParseBool(rwContext.QueryParam("desc"))

	if err != nil {
		desc = false
	}

	if sortType == "" {
		sortType = "flat"
	}

	posts , err := Thread.threadLogic.GetPosts(slugOrId,limit,since,sortType,desc)

	if err != nil {
		return rwContext.JSON(http.StatusNotFound, models.Error{Message:"can't get thread by slug_or_id:" + slugOrId})
	}

	return rwContext.JSON(http.StatusOK,posts)
}

func (Thread ThreadDelivery) UpdateThread(rwContext echo.Context) error {
	slugOrId := rwContext.Param("slug_or_id")
	newThread := new(models.Thread)
	rwContext.Bind(newThread)

	thread , err := Thread.threadLogic.UpdateThread(slugOrId,*newThread)

	if err != nil {
		return rwContext.JSON(http.StatusNotFound, models.Error{"Can't find thread by slug_or_id: "+ slugOrId})
	}

	return rwContext.JSON(http.StatusOK,thread)
}

func (ForumD ThreadDelivery) SetupHandlers(server *echo.Echo) {
	server.POST("/thread/:slug_or_id/create", ForumD.CreatePosts)
	server.POST("/thread/:slug_or_id/vote", ForumD.VoteThread)
	server.POST("/thread/:slug_or_id/details", ForumD.UpdateThread)
	server.GET("/thread/:slug_or_id/details", ForumD.GetThread)
	server.GET("/thread/:slug_or_id/posts", ForumD.GetPosts)

}
