package delivery

import (
	"github.com/labstack/echo"
	"main/internal/models"
	"main/internal/posts"
	"main/internal/posts/usecase"
	"net/http"
	"strconv"
	"strings"
)

type PostDelivery struct {
	postUseLogic posts.PostsUse
}

func NewUserDelivery(postU usecase.PostUseRealisation) PostDelivery {
	return PostDelivery{postUseLogic: postU}
}

func (PostD PostDelivery) GetPost(rwContext echo.Context) error {
	id, _ := strconv.Atoi(rwContext.Param("id"))
	related := rwContext.QueryParams()

	val := strings.Split(related.Get("related"), ",")

	allPostData, err := PostD.postUseLogic.GetPostData(id, val)

	if err != nil {

		return rwContext.JSON(http.StatusNotFound, models.Error{Message: "can't find post by id: " + rwContext.Param("id")})
	}

	return rwContext.JSON(http.StatusOK, allPostData)
}

func (PostD PostDelivery) UpdatePost(rwContext echo.Context) error {
	id, _ := strconv.ParseInt(rwContext.Param("id"), 10, 64)

	msg := new(models.Message)
	rwContext.Bind(msg)

	currentMsg, err := PostD.postUseLogic.UpdatePost(id, msg.Message)

	if err != nil {

		return rwContext.JSON(http.StatusNotFound, models.Error{Message: "can't find post by id: " + rwContext.Param("id")})
	}

	return rwContext.JSON(http.StatusOK, currentMsg)
}

func (PostD PostDelivery) SetupHandlers(server *echo.Echo) {
	server.GET("/post/:id/details", PostD.GetPost)
	server.POST("/post/:id/details", PostD.UpdatePost)
}
