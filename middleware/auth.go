package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/user/daily-records-backend/utils"
)

// AuthMiddleware 验证 Supabase JWT 鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, 401, "未提供认证令牌")
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Error(c, 401, "认证格式不正确")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Supabase JWT 验证逻辑
		// 注意：Supabase JWT 是由其 Auth 服务签发的
		// 验证时需要使用其提供的 JWT Secret (SUPABASE_KEY 在某些场景下可用于简单校验，
		// 但标准解析应通过其 JWK 或共享密钥)
		// 这里我们主要解析出 user_id 并假设外部网关或 Supabase 公钥已配置

		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			utils.Error(c, 401, "无效的令牌")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.Error(c, 401, "解析令牌失败")
			c.Abort()
			return
		}

		// 从 claims 中提取 sub (即 user_id)
		userID, ok := claims["sub"].(string)
		if !ok {
			utils.Error(c, 401, "令牌中不含用户身份信息")
			c.Abort()
			return
		}

		// 将 user_id 注入 Gin 上下文，供后续接口使用
		c.Set("user_id", userID)
		c.Next()
	}
}

// 补充：生产环境下应使用库验证签名，此处为演示提取逻辑
// func verifyToken(tokenString string) bool { ... }
