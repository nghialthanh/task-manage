package routes

import (
	controller "task-manage/controllers"
	"task-manage/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.CORSMiddleware())
	incomingRoutes.POST("/auth/signup", controller.SignUp())
	incomingRoutes.POST("/auth/login", controller.Login())
	incomingRoutes.POST("/auth/refresh-token", controller.RefreshToken())
}
