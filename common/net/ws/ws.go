package ws

//
//import (
//	"fmt"
//	"gameserver/common/logger"
//	"github.com/aceld/zinx/ziface"
//	"github.com/aceld/zinx/zlog"
//	"github.com/aceld/zinx/znet"
//	"net/http"
//)
//
//var S *Server
//
//type Server struct {
//	ziface.IServer
//}
//
//// PingRouter MsgId=1的路由
//type PingRouter struct {
//	znet.BaseRouter
//}
//
//// Ping Handle MsgId=1的路由处理方法
//func (r *PingRouter) Handle(request ziface.IRequest) {
//	//读取客户端的数据
//	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
//}
//
//func auth(r *http.Request) error {
//	token := r.Header.Get("Authorization")
//	if token == "" {
//		return fmt.Errorf("Missing token")
//	}
//	return nil
//}
//
//func NewWSServer() {
//	zlog.SetLogger(new(logger.ZapLogger))
//	s := znet.NewServer()
//	s.SetWebsocketAuth(auth)
//	s.Serve()
//	S = &Server{s}
//}
