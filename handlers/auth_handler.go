package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/user/daily-records-backend/utils"
)

// SignUpRequest 注册请求
type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SignUp 注册接口
func SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "邮箱格式不正确或密码太短")
		return
	}

	// 使用 types.SignupRequest
	res, err := utils.Client.Auth.Signup(types.SignupRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		utils.Error(c, 400, "注册失败: "+err.Error())
		return
	}

	utils.Success(c, res.User)
}

// Login 登录接口
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "请输入正确的邮箱和密码")
		return
	}

	// 使用 SignInWithEmailPassword 直接传参
	res, err := utils.Client.Auth.SignInWithEmailPassword(req.Email, req.Password)
	if err != nil {
		utils.Error(c, 401, "登录失败: 邮箱或密码错误")
		return
	}

	utils.Success(c, gin.H{
		"access_token": res.AccessToken,
		"user":         res.User,
	})
}
