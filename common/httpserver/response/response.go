package response

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// WriteJsonResponse 是一个封装好的函数，用于向客户端返回JSON格式的数据
func WriteJsonResponse(c *gin.Context, data interface{}, code int, message string) {
	resp := &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	c.JSON(code, resp)
}
