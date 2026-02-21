package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Response 统一响应格式
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	userID := c.GetString("user_id")
	path := c.Request.URL.Path

	GetLogger().Info("Request Success",
		zap.String("user_id", userID),
		zap.String("path", path),
		zap.Int("code", 200),
	)

	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	userID := c.GetString("user_id")
	path := c.Request.URL.Path

	GetLogger().Error("Request Error",
		zap.String("user_id", userID),
		zap.String("path", path),
		zap.Int("code", code),
		zap.String("msg", msg),
	)

	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// ValidationError 校验失败响应
func ValidationError(c *gin.Context, msg string) {
	Error(c, 400, msg)
}
