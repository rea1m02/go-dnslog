package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"

	"github.com/rea1m/go-dnslog/database"
	"github.com/rea1m/go-dnslog/models"
	"github.com/spf13/viper"
)

// ListWatchDomains 获取监听域名列表
func ListWatchDomains(c *gin.Context) {
	var watchDomains []models.WatchDomain
	if err := database.DB.Order("created_at DESC").Find(&watchDomains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get watch domains"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"watch_domains": watchDomains})
}

// AddWatchDomain 添加监听域名
func AddWatchDomain(c *gin.Context) {
	var req struct {
		Domain string `json:"domain" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	// 格式化域名（去除末尾的点，转小写）
	domain := strings.TrimSuffix(req.Domain, ".")
	domain = strings.ToLower(domain)

	// 验证域名是否包含主域名后缀
	baseDomain := viper.GetString("dns.domain")
	if !strings.HasSuffix(domain, "."+baseDomain) && domain != baseDomain {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Domain must be a subdomain of " + baseDomain})
		return
	}

	// 检查是否已存在
	var existing models.WatchDomain
	if err := database.DB.Where("domain = ?", domain).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Domain already exists"})
		return
	}

	// 创建监听域名
	watchDomain := models.WatchDomain{
		Domain: domain,
	}

	if err := database.DB.Create(&watchDomain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add watch domain"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"watch_domain": watchDomain})
}

// DeleteWatchDomain 删除监听域名
func DeleteWatchDomain(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	// 查找并删除
	var watchDomain models.WatchDomain
	if err := database.DB.First(&watchDomain, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Watch domain not found"})
		return
	}

	if err := database.DB.Delete(&watchDomain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete watch domain"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Watch domain deleted successfully"})
}
