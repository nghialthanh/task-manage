package middleware

import (
	"net/http"

	helper "task-manage/helpers"

	"github.com/gin-gonic/gin"
)

// Authz validates token and authorizes users
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			helper.SendResponse(c, helper.Response{Status: http.StatusInternalServerError, Message: []string{"No Authorization header provided"}})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			helper.SendResponse(c, helper.Response{Status: http.StatusInternalServerError, Error: []string{err}})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("user_type", claims.User_type)
		c.Set("user_id", claims.Uid)

		c.Next()

	}
}
