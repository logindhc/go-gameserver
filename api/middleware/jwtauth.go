package middleware

import (
	"gameserver/api/url"
	"gameserver/common/utils"
	"gameserver/core/redis"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == url.USER_LOGIN_URL {
			c.Next()
			return
		}
		// 从header中获取token
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Missing token",
			})
			c.Abort()
			return
		}
		split := strings.Split(token, " ")
		if len(split) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token format",
			})
			c.Abort()
			return
		}
		token = split[1]
		// 验证token
		openId, err := utils.GetJwtOpenId(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}
		userId := redis.GetLoginToken(openId)
		if userId == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}
		c.Set("token", userId)
		c.Next()
	}
}
