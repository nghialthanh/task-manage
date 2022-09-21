package main

import (
	"fmt"
	"os"
	"time"

	routes "task-manage/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	str := time.Now().Local().Add(time.Hour * time.Duration(24))
	fmt.Println(time.Now().Local())
	fmt.Println(str)
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.ProjectRoutes(router)

	router.Run(":" + port)
}
