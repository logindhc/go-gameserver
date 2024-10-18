package utils

import (
	"encoding/base64"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
	"time"
)

var (
	JWT_SECRET  = []byte("^JWT#key#@byonegames.com")
	JWT_ISS     = "byonegames"
	JWT_EXP     = int64(3600 * 24)
	JWT_SUBJECT = "token"
)

func CreateJwt(openId string, expTime int64) (string, error) {
	now := time.Now().Unix()
	exp := now + expTime
	if expTime == 0 {
		exp = now + JWT_EXP
	}
	nowStr := strconv.FormatInt(now, 10)
	openId += "_"
	openId += nowStr
	encryptCBC := cryptor.AesCbcEncrypt([]byte(openId), JWT_SECRET)
	encryptOpenId := base64.StdEncoding.EncodeToString(encryptCBC)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		JWT_SUBJECT: encryptOpenId,
		"exp":       exp,
		"iss":       JWT_ISS,
		"iat":       now,
	})
	return claims.SignedString(JWT_SECRET)
}

func GetJwtOpenId(jwtToken string) (string, error) {
	mapClaims, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) { return JWT_SECRET, nil })
	if err != nil || !mapClaims.Valid {
		return "", err
	}
	if claims, ok := mapClaims.Claims.(jwt.MapClaims); ok {
		for key, val := range claims {
			if key == JWT_SUBJECT {
				decodeString, err := base64.StdEncoding.DecodeString(val.(string))
				if err != nil {
					return "", err
				}
				openId := string(cryptor.AesCbcDecrypt(decodeString, JWT_SECRET))
				openId = strings.Split(openId, "_")[0]
				return openId, nil
			}
		}

	}
	return "", err
}
