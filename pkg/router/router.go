package router

import (
	"chat/v1/pkg/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter(ctrl controller.Controller) *gin.Engine {
	router := gin.Default()

	router.GET("/channel/:id/:user/stream/view", ctrl.Channel().StreamHandler)

	router.POST("/channel/:id/message/send", ctrl.Channel().SendMassageHandler)

	return router
}
