package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	commonLogPath string
	errorLogPath  string
	commonFd      *os.File
	errorFd       *os.File
)

func SetupGinLog() {
	// 获取当前工作路径
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get current working directory")
	}

	// 如果 *LogDir 为空，则使用工作路径下的默认日志目录
	if *LogDir == "" {
		*LogDir = filepath.Join(cwd, "logs")
	}

	// 创建日志目录（如果不存在）
	err = os.MkdirAll(*LogDir, 0755)
	if err != nil {
		log.Fatal("failed to create log directory")
	}

	// 每天凌晨1点执行日志清理任务
	go func() {
		for {
			next := time.Now().Add(24 * time.Hour)
			next = time.Date(next.Year(), next.Month(), next.Day(), 1, 0, 0, 0, next.Location())
			timer := time.NewTimer(next.Sub(time.Now()))

			<-timer.C
			cleanupLogs(*LogDir)
			initLogFiles()
		}
	}()

	// 初始化日志文件描述符
	initLogFiles()

	// 设置 Gin 的日志输出
	gin.DefaultWriter = io.MultiWriter(os.Stdout, commonFd)
	gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, errorFd)
}

func initLogFiles() {
	// 获取当前日期
	date := time.Now().Format("2006-01-02")

	// 构建日志文件路径
	commonLogPath = filepath.Join(*LogDir, fmt.Sprintf("common-%s.log", date))
	errorLogPath = filepath.Join(*LogDir, fmt.Sprintf("error-%s.log", date))

	// 关闭旧的日志文件描述符
	if commonFd != nil {
		commonFd.Close()
	}
	if errorFd != nil {
		errorFd.Close()
	}

	// 打开新的日志文件描述符
	var err error
	commonFd, err = os.OpenFile(commonLogPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("failed to open common log file")
	}

	errorFd, err = os.OpenFile(errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("failed to open error log file")
	}
}

func cleanupLogs(logDir string) {
	// 列出日志目录下的所有文件
	files, err := filepath.Glob(filepath.Join(logDir, "*.log"))
	if err != nil {
		log.Println("failed to list log files:", err)
		return
	}

	// 遍历所有日志文件
	for _, file := range files {
		// 解析日志文件名中的日期
		base := filepath.Base(file)
		dateStr := base[len("common-") : len(base)-len(".log")]
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Println("failed to parse log file date:", err)
			continue
		}

		// 如果日志文件日期早于3天前，则删除该文件
		if date.Before(time.Now().AddDate(0, 0, -3)) {
			err := os.Remove(file)
			if err != nil {
				log.Println("failed to remove log file:", err)
			}
		}
	}
}

func SysLog(s string) {
	t := time.Now()
	_, _ = fmt.Fprintf(gin.DefaultWriter, "[SYS] %v | %s \n", t.Format("2006/01/02 - 15:04:05"), s)
}

func SysError(s string) {
	t := time.Now()
	_, _ = fmt.Fprintf(gin.DefaultErrorWriter, "[SYS] %v | %s \n", t.Format("2006/01/02 - 15:04:05"), s)
}

func FatalLog(v ...any) {
	t := time.Now()
	_, _ = fmt.Fprintf(gin.DefaultErrorWriter, "[FATAL] %v | %v \n", t.Format("2006/01/02 - 15:04:05"), v)
	os.Exit(1)
}
