package routes

import (
	handler2 "gameserver/appserver/api/handler"
	"gameserver/appserver/url"
	"github.com/gin-gonic/gin"
)

func New(engine *gin.Engine) {
	engine.POST(url.USER_LOGIN_URL, handler2.Login)
	engine.POST(url.USER_INFO_URL, handler2.GetUserInfo)
	engine.POST(url.USER_LEVEL_URL, handler2.UpLevel)
	engine.POST(url.USER_UPDATE_URL, handler2.Update)
}
