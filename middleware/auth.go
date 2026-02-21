package middleware

import (
	"os"
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
		var claims jwt.MapClaims
		var token *jwt.Token
		var err error

		jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")
		if jwtSecret != "" {
			// 1. 如果配置了密钥，进行完整签名验证 (生产环境推荐)
			token, err = jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
		} else {
			// 2. 如果未配置密钥，仅解析声明 (仅供演示或信任前端安全场景)
			token, _, err = new(jwt.Parser).ParseUnverified(tokenString, &jwt.MapClaims{})
		}

		if err != nil || token == nil {
			utils.Error(c, 401, "无效或已过期的令牌")
			c.Abort()
			return
		}

		// 解析声明
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			// 如果是 ParseWithClaims 可能会返回 *jwt.MapClaims
			if mapPtr, ok := token.Claims.(*jwt.MapClaims); ok {
				claims = *mapPtr
			} else {
				utils.Error(c, 401, "解析令牌声明失败")
				c.Abort()
				return
			}
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
