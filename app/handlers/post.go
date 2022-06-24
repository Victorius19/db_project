package handlers

import (
	"db_project/app/models"
	"db_project/app/usecases"
	"db_project/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"
	"strings"
)

type HandlerPosts struct {
	UseCase usecases.IPostUseCase
}

func MakePostsHandler(useCase usecases.IPostUseCase) *HandlerPosts {
	return &HandlerPosts{UseCase: useCase}
}

func (handler *HandlerPosts) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest)
		return
	}

	detailsRaw := c.Query("related")
	var details []string
	if detailsRaw != "" {
		details = strings.Split(detailsRaw, ",")
	}

	post, err := handler.UseCase.Get(int(id), details)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, post)
	return
}

func (handler *HandlerPosts) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest)
		return
	}

	post := &models.Post{}
	err = easyjson.UnmarshalFromReader(c.Request.Body, post)
	if err != nil {
		c.AbortWithStatusJSON(errors.BadRequest.Code(), errors.BadRequest)
		return
	}

	post.ID = int(id)

	forum, err := handler.UseCase.Update(post)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, forum)
	return
}
