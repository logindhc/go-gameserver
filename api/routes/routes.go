package routes

import (
	"gameserver/api/handler"
	"gameserver/api/url"
	"github.com/gin-gonic/gin"
)

func New(engine *gin.Engine) {
	engine.POST(url.USER_LOGIN_URL, handler.Login)
	engine.POST(url.USER_INFO_URL, handler.GetUserInfo)
	engine.POST(url.USER_LEVEL_URL, handler.UpLevel)
	engine.POST(url.USER_UPDATE_URL, handler.Update)
}
