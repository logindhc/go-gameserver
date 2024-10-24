package middleware

import (
	"fmt"
	"gameserver/core/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	// 获取WebSocket连接
	wsUpgrader = websocket.Upgrader{
		HandshakeTimeout: time.Second * 10,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			//token := r.Header.Get("Authorization")
			//if token == "" {
			//	fmt.Println("token is null")
			//	return false
			//}
			return true
		},
	}
)

func WsHandler(c *gin.Context) {
	if websocket.IsWebSocketUpgrade(c.Request) {
		ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		go func() {
			// 处理WebSocket消息
			for {
				messageType, p, err := ws.ReadMessage()
				if err != nil {
					logger.Logger.Info(err.Error())
					return
				}
				switch messageType {
				case websocket.TextMessage:
					fmt.Printf("处理文本消息, %s\n", string(p))
					ws.WriteMessage(websocket.TextMessage, p)
					// c.Writer.Write(p)
				case websocket.BinaryMessage:
					fmt.Println("处理二进制消息")
				case websocket.CloseMessage:
					fmt.Println("关闭websocket连接")
					return
				case websocket.PingMessage:
					fmt.Println("处理ping消息")
					ws.WriteMessage(websocket.PongMessage, []byte("ping"))
				case websocket.PongMessage:
					fmt.Println("处理pong消息")
					ws.WriteMessage(websocket.PongMessage, []byte("pong"))
				default:
					fmt.Printf("未知消息类型: %d\n", messageType)
					return
				}
			}
		}()
	} else {
		logger.Logger.Error("非websocket请求")
		//处理普通请求
		c.Next()
	}

}
