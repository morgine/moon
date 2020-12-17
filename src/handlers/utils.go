package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/morgine/moon/src/errors"
	"net/http"
)

type Message struct {
	Status  errors.Code
	Message string
	Data    interface{}
}

func SendMessage(ctx *gin.Context, code errors.Code, message string) {
	ctx.AbortWithStatusJSON(http.StatusOK, Message{
		Status:  code,
		Message: message,
	})
}

func SendJSON(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, Message{
		Status: errors.StatusOK,
		Data:   data,
	})
}

func SendError(ctx *gin.Context, err error) {
	code, ok := errors.Unwrap(err)

	if ok {
		ctx.AbortWithStatusJSON(http.StatusOK, Message{
			Status:  code,
			Message: errors.Texts[code],
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusOK, Message{
			Status:  errors.StatusUnknown,
			Message: err.Error(),
		})
	}
}
