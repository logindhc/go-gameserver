package middleware

import (
	"fmt"
	"gameserver/api/handler"
	"gameserver/common/net/http/response"
	"gameserver/common/utils"
	"gameserver/conf"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MD5Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path != "/login" {
			c.Next()
			return
		}
		loginBody := new(handler.LoginBody)
		if err := c.ShouldBindJSON(&loginBody); err != nil {
			response.WriteJsonResponse(c, nil, http.StatusBadRequest, "invalid request param")
			c.Abort()
			return
		}
		time := loginBody.Time
		sign := fmt.Sprintf("%s%d", conf.GameConfig.App.MisKey, time)
		if !utils.Sign(loginBody.Sign, sign) {
			response.WriteJsonResponse(c, nil, http.StatusBadRequest, "sign error")
			c.Abort()
			return
		}
		c.Next()
	}
}
