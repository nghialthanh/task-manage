package routes

import (
	controller "task-manage/controllers"
	"task-manage/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.CORSMiddleware())
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/user/list-users", controller.GetListUsers())
	incomingRoutes.GET("/user/info-users", controller.GetUserByToken())
	incomingRoutes.GET("/users/:user_id", controller.GetUserByID()) // remove
}
