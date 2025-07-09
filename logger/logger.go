package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"webhook/config"
)

// NewLogger 创建并返回一个新的 lumberjack.Logger 实例，用于日志记录。
// 该函数不接受任何参数。
// 返回值是 *lumberjack.Logger，一个指向 lumberjack.Logger 结构的指针，
// 配置了基本的日志文件路径和日志文件的滚动策略。
func NewLogger() *lumberjack.Logger {
	// 返回一个新的 lumberjack.Logger 实例，配置了日志文件路径和滚动策略。
	return &lumberjack.Logger{
		// Filename: 日志文件的路径，通过调用 config.GetLogFilePath() 获取。
		Filename: config.GetLogFilePath(),
		// MaxSize: 单个日志文件的最大大小为 10MB。
		MaxSize: 10,
		// MaxBackups: 不保留旧的日志文件备份。
		MaxBackups: 0,
		// MaxAge: 不限制日志文件的最长保存时间。
		MaxAge: 0,
		// Compress: 不压缩旧的日志文件备份。
		Compress: false,
	}
}
