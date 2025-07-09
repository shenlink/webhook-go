package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"webhook/config"
	"webhook/gitee"
	"webhook/logger"
)

func main() {
	// 加载环境变量
	config.LoadEnv()
	// 配置日志
	log.SetOutput(logger.NewLogger())
	// 创建路由
	router := gin.Default()

	// 首页
	router.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "shenlink gitee/github webhook")
	})

	// Gitee Webhook 路由分组
	giteeGroup := router.Group("/api/gitee")
	// Gitee blog 仓库
	giteeGroup.POST("/blog", gitee.Blog)

	// 获取端口
	port := config.GetPort()
	// 启动服务
	router.Run(":" + strconv.Itoa(port))
}
