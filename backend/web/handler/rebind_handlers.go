package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"github.com/rea1m/go-dnslog/models"
	"github.com/spf13/viper"
	"github.com/rea1m/go-dnslog/database"
	"crypto/md5"
	"encoding/hex"
)


// RebindList
func RebindList(c *gin.Context) {
	userID, _ := c.Get("userID")
	var rebindList []models.Rebind
	database.DB.Where("user_id = ?", userID).Find(&rebindList)
	c.JSON(http.StatusOK, gin.H{"rebind_list": rebindList})
}


// RebindGen 生成DNS Rebind记录
func RebindGen(c *gin.Context) {
	var req struct {
		FirstIp   	string 	`json:"first_ip" binding:"required"`
		SecondIp 	string 	`json:"second_ip" binding:"required"`
	}

	userID, _ := c.Get("userID")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	dnsDomain := viper.GetString("dns.domain")
	hash := md5Hash(req.FirstIp + req.SecondIp)
	// 生成Rebind域名
	rebindDomain := fmt.Sprintf("%s.e.%s", hash, dnsDomain)

	// // 检查是否存在相同的哈希值
	// var existingRebind models.Rebind
	// if err := database.DB.Where("hash = ? AND user_id = ?", hash, userID).First(&existingRebind).Error; err == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Hash already exists"})
	// 	return
	// }

	// 保存到数据库
	rebind := models.Rebind{
		Domain: rebindDomain,
		UserID:    userID.(uint),
		FirstIP:   req.FirstIp,
		SecondIP:  req.SecondIp,
	}
	if err := database.DB.Create(&rebind).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rebind record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rebind_domain": rebindDomain})
}

func RebindDelete(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	var rebind models.Rebind
	if err := database.DB.Where("id = ? AND user_id = ?", req.ID, userID).First(&rebind).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rebind record not found"})
		return
	}

	if err := database.DB.Delete(&rebind).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rebind record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rebind record deleted successfully"})
}

// md5Hash 计算字符串的MD5哈希值
func md5Hash(str string) string {
	hash := md5.New()
	_, _ = hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil)[:4])
}