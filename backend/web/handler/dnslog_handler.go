package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"

	"github.com/rea1m/go-dnslog/database"
	"github.com/rea1m/go-dnslog/models"
)

// ListDNSLogs 获取DNS日志列表
// 前端设置相应按钮，实现自动刷新
func ListDNSLogs(c *gin.Context) {
	var req struct {
		PageNumber int    `json:"pageNumber" binding:"required,min=1"`
		PageSize   int    `json:"pageSize" binding:"required,min=1,max=100"`
		Search     string `json:"search"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	db := database.DB.Model(&models.DNSLog{}).Where("user_id = ?", userID)

	if req.Search != "" {
		escapedSearch := strings.ReplaceAll(req.Search, "%", "\\%")
		escapedSearch = strings.ReplaceAll(escapedSearch, "_", "\\_")
		db = db.Where("host LIKE ? OR ip LIKE ?", "%"+escapedSearch+"%", "%"+escapedSearch+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var logs []models.DNSLog
	offset := (req.PageNumber - 1) * req.PageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      req.PageNumber,
		"page_size": req.PageSize,
		"logs":      logs,
	})
}

// DeleteDNSLogs 删除单个DNS日志
func DeleteDNSLogs(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ID uint `json:"id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	// 先查询记录是否存在且属于当前用户
	var dnsLog models.DNSLog
	result := database.DB.Where("id = ? AND user_id = ?", req.ID, userID).First(&dnsLog)

	if result.RowsAffected == 0 {
		// 记录不存在或不属于当前用户
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this DNS log"})
		return
	}

	// 执行删除
	database.DB.Delete(&dnsLog)

	c.JSON(http.StatusOK, gin.H{"message": "DNS log deleted successfully"})
}

func BatchDeleteDNSLogs(c *gin.Context) {
	userID, _ := c.Get("userID")

	// 删除所有日志
	database.DB.Where("user_id = ?", userID).Delete(&models.DNSLog{})

	c.JSON(http.StatusOK, gin.H{"message": "DNS logs deleted successfully"})
}
