package gitee

import (
	"log"
	"net/http"
	"webhook/command"
	"webhook/config"

	"github.com/gin-gonic/gin"
)

// 控制最大并发执行 shell 命令的数量
var maxShellExecConcurrentSize = config.GetMaxShellExecConcurrentSize()
var maxShellExecConcurrent = make(chan struct{}, maxShellExecConcurrentSize)

// 处理Gitee blog 仓库的 Webhook请求
// 参数：c 请求上下文
func Blog(c *gin.Context) {
	// 获取请求头中的token和timestamp
	token := c.GetHeader("X-Gitee-Token")
	timestamp := c.GetHeader("X-Gitee-Timestamp")

	// 验证token的合法性
	status, err := validateToken(token, timestamp)
	if err != nil {
		// 如果验证失败，返回错误信息并记录日志
		c.String(status, err.Error())
		log.Printf("验证gitee token失败: %v", err)
		return
	}

	// 返回成功响应
	c.String(http.StatusOK, "webhook处理成功")

	// 执行shell命令更新博客内容
	command.ExecuteShellCommandAsync("cd /www/wwwroot/blog && git pull origin master && npm run docs:build", maxShellExecConcurrent)
}
