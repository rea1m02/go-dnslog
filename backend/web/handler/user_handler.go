package handler

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/rea1m/go-dnslog/database"
	"github.com/rea1m/go-dnslog/models"
	"github.com/rea1m/go-dnslog/web/middleware"
)

// Index 处理根路径请求
// 无需处理
//func Index(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{"message": "Welcome to go-dnslog API"})
//}

// Login 用户登录
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	// 查询用户
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username or password is incorrect"})
		return
	}

	// 检查登录失败次数
	if user.TryLoginCounter > 10 && time.Since(user.LastTryLoginTime) < 10*time.Minute {
		c.JSON(http.StatusForbidden, gin.H{"error": "Too many login attempts, please try again later"})
		return
	}

	// 验证密码
	salt := viper.GetString("security.password_salt")
	passwordHash := md5.Sum([]byte(req.Password + req.Username[:3] + salt))
	passwordHex := hex.EncodeToString(passwordHash[:])

	if user.Password != passwordHex {
		user.TryLoginCounter++
		user.LastTryLoginTime = time.Now()
		database.DB.Save(&user)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username or password is incorrect"})
		return
	}

	// 生成JWT令牌
	jwtToken, err := middleware.GenerateToken(user.ID, user.Username, user.JWTTokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 更新用户登录信息
	user.TryLoginCounter = 0
	user.LoginIP = c.ClientIP()
	database.DB.Save(&user)

	host := viper.GetString("dns.domain")

	c.JSON(http.StatusOK, gin.H{
		"token":       jwtToken,
		"username":    user.Username,
		"user_domain": user.UserDomain,
		"is_admin":    user.IsAdmin,
		"host":        host,
	})

}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":    user.Username,
		"email":       user.Email,
		"user_domain": user.UserDomain,
		"token":       user.Token,
		"is_admin":    user.IsAdmin,
	})
}

// Logout 用户登出
// 由于使用jwt，所以登出功能需要依赖外部存储，如redis，这里不做考虑， 仅在前端通过清除token以及返回登录页面的方式实现
//func Logout(c *gin.Context) {
//	userID, exists := c.Get("userID")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
//		return
//	}
//
//	// 递增用户JWT令牌版本
//	var user models.User
//	if err := database.DB.First(&user, userID).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
//		return
//	}
//
//	user.JWTTokenVersion++
//	if err := database.DB.Save(&user).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
//}
