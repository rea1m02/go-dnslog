package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"

	"github.com/rea1m/go-dnslog/database"
	"github.com/rea1m/go-dnslog/models"
)

// JWTClaims 定义JWT载荷结构
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Username string `json:"username"`
	TokenVersion uint `json:"token_version"`
	jwt.RegisteredClaims
}

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// 解析JWT令牌
		secret := viper.GetString("security.jwt_secret")
		claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 验证用户是否存在
	var user models.User
	if err := database.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		c.Abort()
		return
	}

	// 验证令牌版本
	if claims.TokenVersion != user.JWTTokenVersion {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token version"})
		c.Abort()
		return
	}

	// 将用户信息存储到上下文
		c.Set("userID", user.ID)
		c.Set("username", user.Username)
		c.Set("isAdmin", user.IsAdmin)

		c.Next()
	}
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string, tokenVersion uint) (string, error) {
	// 设置过期时间
	expiry := viper.GetInt64("security.token_expiry")
	claims := JWTClaims{
	UserID:   userID,
	Username: username,
	TokenVersion: tokenVersion,
	RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-dnslog",
		},
	}

	// 创建令牌
	secret := viper.GetString("security.jwt_secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err

}