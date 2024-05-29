package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"wechat-server/common"

	"github.com/gin-gonic/gin"
)

// GetCommonLogs handles GET requests to fetch common logs of the current day
func GetCommonLogs(c *gin.Context) {
	logs, err := readLogs("common")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to read common logs: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "",
		"success": true,
		"data":    logs,
	})
}

// GetErrorLogs handles GET requests to fetch error logs of the current day
func GetErrorLogs(c *gin.Context) {
	logs, err := readLogs("error")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to read error logs: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "",
		"success": true,
		"data":    logs,
	})
}

// readLogs reads logs from the file specified by log type
func readLogs(logType string) (string, error) {
	if *common.LogDir == "" {
		return "", fmt.Errorf("log directory path is not set")
	}

	today := time.Now().Format("2006-01-02")             // 获取当前日期
	fileName := fmt.Sprintf("%s-%s.log", logType, today) // 构造文件名
	filePath := filepath.Join(*common.LogDir, fileName)  // 构建完整的文件路径

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("log file does not exist: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading log file: %s, error: %v", filePath, err)
	}
	return string(content), nil
}
