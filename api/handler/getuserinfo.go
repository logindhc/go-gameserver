package handler

import (
	"gameserver/common/net/http/response"
	"gameserver/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(ctx *gin.Context) {
	user := models.UserRepository.Get(ctx.GetInt64("token"))
	if user == nil {
		response.WriteJsonResponse(ctx, nil, http.StatusInternalServerError, "error")
		return
	}
	response.WriteJsonResponse(ctx, user, http.StatusOK, "ok")
}
