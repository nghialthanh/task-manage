package routes

import (
	controller "task-manage/controllers"

	"github.com/gin-gonic/gin"
)

//UserRoutes function
func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/auth/signup", controller.SignUp())
	incomingRoutes.POST("/auth/login", controller.Login())
	incomingRoutes.POST("/auth/refresh-token", controller.RefreshToken())
}
