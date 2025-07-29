package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


// NotFound 处理404
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
}
