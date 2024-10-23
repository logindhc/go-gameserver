package handler

import (
	"gameserver/appserver/service"
	"gameserver/common/net/http/response"
	"gameserver/common/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginBody struct {
	Code     string `json:"code"`
	Platform int    `json:"platform"`
	Channel  int    `json:"channel"`
	Time     int64  `json:"time"`
	Sign     string `json:"sign"`
}

func Login(ctx *gin.Context) {
	loginBody := new(LoginBody)
	if err := ctx.ShouldBindJSON(&loginBody); err != nil {
		response.WriteJsonResponse(ctx, nil, http.StatusInternalServerError, "error")
		return
	}
	account, err2 := service.GetOrCreateAccount(loginBody.Code, loginBody.Platform, loginBody.Channel)
	if err2 != nil {
		response.WriteJsonResponse(ctx, nil, http.StatusInternalServerError, "error")
		return
	}
	user, err2 := service.Login(account)
	if err2 != nil {
		response.WriteJsonResponse(ctx, nil, http.StatusInternalServerError, "error")
		return
	}

	jwt, err := utils.CreateJwt(user.OpenId, 0)
	if err != nil {
		response.WriteJsonResponse(ctx, nil, http.StatusInternalServerError, "error")
		return
	}
	ctx.Header("Authorization", "Bearer "+jwt)

	response.WriteJsonResponse(ctx, account, http.StatusOK, "ok")
}
