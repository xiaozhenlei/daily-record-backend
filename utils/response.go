package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应格式
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{ // 要求统一返回 200 或业务 code，这里按需求返回 JSON
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// ValidationError 校验失败响应
func ValidationError(c *gin.Context, msg string) {
	Error(c, 400, msg)
}
