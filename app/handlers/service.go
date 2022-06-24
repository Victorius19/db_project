package handlers

import (
	"db_project/app/usecases"
	"db_project/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerServices struct {
	UseCase usecases.IServiceUseCase
}

func MakeServicesHandler(useCase usecases.IServiceUseCase) *HandlerServices {
	return &HandlerServices{UseCase: useCase}
}

func (handler *HandlerServices) Clear(c *gin.Context) {
	err := handler.UseCase.Clear()
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (handler *HandlerServices) Status(c *gin.Context) {
	status, err := handler.UseCase.Status()
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.MsgErrors).Code(), err)
		return
	}

	c.JSON(http.StatusOK, status)
}
