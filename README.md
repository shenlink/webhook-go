# Webhook 服务说明

这是一个基于 Go 编写的 Webhook 接收服务，主要用于接收 Gitee 的 Webhook 请求并执行相应的自动化操作（如自动拉取代码、构建文档等）。

## 功能特点

- 验证 Gitee Webhook 请求的签名和时间戳，确保请求来源合法。
- 支持异步执行 Shell 命令，避免阻塞主线程。
- 可配置最大并发执行数量，防止资源过载。
- 日志输出支持文件滚动切割，便于日志管理。

## 技术栈

- Go (Golang)
- [Gin](https://github.com/gin-gonic/gin) 框架
- [Lumberjack](https://github.com/natefinch/lumberjack.v2) 日志切割库
- 使用 HMAC-SHA256 进行安全验证

## 目录结构

```
.
├── config/          # 配置模块
│   └── config.go    # 配置读取与环境变量处理
├── main.go          # 主程序入口
├── README.md        # 项目说明文档
└── LICENSE          # 开源协议文件
```

## 环境变量配置

请在 `.env` 文件中配置以下参数：

| 参数名 | 含义 | 默认值 |
|--------|------|--------|
| `TIMESTAMP_TOLERANCE` | 允许的时间戳误差（秒） | `300` |
| `GITEE_SIGN_KEY` | Gitee Webhook 的密钥 | 无默认值 |
| `MAX_CONCURRENT` | 最大并发执行数 | `1` |
| `PORT` | 服务监听端口 | `8080` |
| `LOG_FILE_PATH` | 日志文件路径 | `logs/app.log` |

## 快速启动

1. 安装依赖：

```bash
go mod tidy
```

2. 启动服务：

```bash
go run main.go
```

3. 访问首页测试接口：

```
GET http://localhost:8080
```

4. 触发 Gitee Webhook：

```
POST http://localhost:8080/api/gitee/blog
```

## 接口说明

### 首页

- 方法：GET
- 路径：`/`
- 返回：简单的欢迎信息

### Gitee Webhook 接口

- 方法：POST
- 路径：`/api/gitee/blog`
- 请求头：
  - `X-Gitee-Token`: Gitee 提供的签名值
  - `X-Gitee-Timestamp`: 时间戳（毫秒）
- 行为：验证签名后执行指定的 Shell 命令（如 `git pull` 和 `npm run docs:build`）

## 安全提示

- 请务必设置 `GITEE_SIGN_KEY`，否则无法验证请求合法性。
- 不建议将服务直接暴露在公网，请通过反向代理或防火墙保护服务。

## 开源协议

本项目使用 MIT 协议开源，请参见 [LICENSE](./LICENSE) 文件。

## 贡献者

欢迎提交 PR 或 Issue！