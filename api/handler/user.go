package handler

import (
	"gameserver/common/net/http/response"
	"gameserver/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpLevel(ctx *gin.Context) {
	type ls struct {
		Level int `json:"level"`
	}
	level := new(ls)
	err := ctx.BindJSON(&level)
	if err != nil {
		response.WriteJsonResponse(ctx, nil, http.StatusBadRequest, "bad request")
		return
	}
	user := models.UserRepository.Get(ctx.GetInt64("token"))
	user.Level = level.Level

	models.UserRepository.Update(user)

	response.WriteJsonResponse(ctx, user, http.StatusOK, "ok")
}

func Update(ctx *gin.Context) {
	type Item struct {
		Id    int `json:"id"`
		Count int `json:"count"`
	}
	item := new(Item)
	err := ctx.BindJSON(&item)
	if err != nil {
		response.WriteJsonResponse(ctx, nil, http.StatusBadRequest, "bad request")
		return
	}
	user := models.UserRepository.Get(ctx.GetInt64("token"))
	user.GetItems()[item.Id] = item.Count
	models.UserRepository.Update(user)
	models.UserRepository.Flush()

	response.WriteJsonResponse(ctx, user, http.StatusOK, "ok")
}
