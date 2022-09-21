package helper

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int
	Message []string
	Error   []string
	Data    interface{}
	Code    string
}

func SendResponse(c *gin.Context, response Response) {
	if len(response.Message) > 0 {
		c.JSON(response.Status, map[string]interface{}{"status": response.Status, "message": strings.Join(response.Message, "; ")})
	} else if len(response.Error) > 0 {
		c.JSON(response.Status, map[string]interface{}{"status": response.Status, "error": strings.Join(response.Error, "; ")})
	} else {
		c.JSON(response.Status, map[string]interface{}{"status": response.Status, "data": response.Data})
	}
}
