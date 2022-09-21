package routes

import (
	controller "task-manage/controllers"
	"task-manage/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func ProjectRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.POST("/project/create", controller.CreateProject())
}
