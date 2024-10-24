package service

import (
	"fmt"
	"gameserver/common/utils"
	"gameserver/core/redis"
	"gameserver/models"
	"time"
)

/*
获取账号信息
*/
func GetOrCreateAccount(openId string, platform int, channel int) (*models.Account, error) {
	// 获取数据库操作对象
	// 直接获取
	accountId := fmt.Sprintf("%d_%s", channel, openId)
	account := models.AccountRepository.Get(accountId)
	if account == nil {
		account = &models.Account{
			ID:            accountId,
			OpenId:        utils.UUID2STR(),
			Channel:       channel,
			ChannelOpenId: openId,
			Platform:      platform,
			RegisterTime:  time.Now().Unix(),
		}
		models.AccountRepository.Add(account)
	}
	return account, nil
}

func Login(account *models.Account) (*models.User, error) {
	userOpenId := account.OpenId
	// 从redis 根据游戏生成的openId获取userId
	userId := redis.GetLoginToken(userOpenId)
	user := models.User{}
	if userId == 0 {
		//where 查询
		models.UserRepository.Where("open_id", userOpenId).Find(&user)
		if user.ID != 0 {
			userId = user.ID
		}
	}
	if userId == 0 { //新玩家
		user = models.User{
			ID:            utils.UUID2INT(),
			OpenId:        userOpenId,
			Channel:       account.Channel,
			Platform:      account.Platform,
			RegisterTime:  account.RegisterTime,
			LastLoginTime: time.Now().Unix(),
		}
		models.UserRepository.Add(&user)
	}
	redis.SetLoginToken(userOpenId, user.ID)
	return &user, nil
}
