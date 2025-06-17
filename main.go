package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"

	"log"
	"os/exec"
	"strconv"
	"time"
	"webhook/config"
)

// 控制最大并发执行 shell 命令的数量
var maxShellExecConcurrent chan struct{}

// 处理Gitee blog 仓库的 Webhook请求
func giteeBlog(c *gin.Context) {
	err := validateGiteeToken(c)
	if err != nil {
		log.Printf("验证gitee token失败: %v", err)
		return
	}

	// 执行shell命令
	ExecuteShellCommandAsync("cd /www/wwwroot/blog && git pull origin master && npm run docs:build")

	// 返回成功响应
	c.String(http.StatusOK, "webhook处理成功")
}

func validateGiteeToken(c *gin.Context) error {
	// 获取Gitee请求头
	giteeToken := c.GetHeader("X-Gitee-Token")
	giteeTimestamp := c.GetHeader("X-Gitee-Timestamp")

	if giteeToken == "" || giteeTimestamp == "" {
		c.String(http.StatusBadRequest, "没有header头")
		return errors.New("没有header头")
	}

	// 从配置中读取时间差
	timestampTolerance := config.GetTimestampTolerance()

	// 验证时间戳是否在有效期内
	currentTime := time.Now().Unix()
	timestamp, err := strconv.ParseInt(giteeTimestamp, 10, 64)
	if err != nil {
		c.String(http.StatusInternalServerError, "时间戳格式错误")
		return errors.New("时间戳格式错误")
	}

	// 将毫秒转换为秒
	timestampInSeconds := timestamp / 1000

	if abs(currentTime-timestampInSeconds) > int64(timestampTolerance) {
		c.String(http.StatusUnauthorized, "时间戳不正确")
		return errors.New("时间戳不正确")
	}

	// 计算签名
	signKey := config.GetGiteeSignKey()
	secStr := fmt.Sprintf("%s\n%s", giteeTimestamp, signKey)

	mac := hmac.New(sha256.New, []byte(signKey))
	mac.Write([]byte(secStr))
	computeToken := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// 安全比较token
	if !hmac.Equal([]byte(computeToken), []byte(giteeToken)) {
		c.String(http.StatusUnauthorized, "token不正确")
		return errors.New("token不正确")
	}
	return nil
}

// 计算两个数的绝对值
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func ExecuteShellCommandAsync(command string) {

	maxShellExecConcurrent <- struct{}{}

	go func() {
		// 执行 shell 命令
		cmd := exec.Command("sh", "-c", command)
		// 不能缺少 HOME=/root 否则会报错:
		// fatal: detected dubious ownership in repository at '/www/wwwroot/blog'
		// To add an exception for this directory, call:
		// git config --global --add safe.directory /www/wwwroot/blog
		cmd.Env = append(os.Environ(), "HOME=/root")
		// 获取输出结果
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("执行shell命令失败，失败原因: %s", output)
			return
		}
		log.Printf("执行成功，输出结果: %s", output)

		<-maxShellExecConcurrent
	}()
}

// 新增日志配置函数
func newLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   config.GetLogFilePath(),
		MaxSize:    10,
		MaxBackups: 0,
		MaxAge:     0,
		Compress:   false,
	}
}

func main() {
	// 加载环境变量
	config.LoadEnv()
	// 配置日志
	log.SetOutput(newLogger())

	// 从配置中读取执行 shell 命令的最大并发数
	maxShellExecConcurrentSize := config.GetMaxShellExecConcurrentSize()
	maxShellExecConcurrent = make(chan struct{}, maxShellExecConcurrentSize)
	port := config.GetPort()

	router := gin.Default()

	// 首页
	router.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "shenlink gitee/github webhook")
	})

	// Gitee Webhook 路由分组
	giteeGroup := router.Group("/api/gitee")
	giteeGroup.POST("/blog", giteeBlog)

	// 启动服务
	router.Run(":" + strconv.Itoa(port))
}
