package routes

import (
	controller "task-manage/controllers"
	"task-manage/middleware"

	"github.com/gin-gonic/gin"
)

//UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/users", controller.GetListUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUserByID())
}
