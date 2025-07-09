package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv 加载环境变量文件。
// 该函数尝试加载名为.env的环境变量文件，如果文件存在则加载到环境中。
// 如果加载失败，将记录错误信息并停止程序执行。
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("加载 .env 文件失败: %v", err)
	}
}

// getEnv 返回指定环境变量的值。如果环境变量未设置，则返回默认值。
// 这个函数的目的是提供一种简便的方式来获取环境变量的值，并在变量未设置时有一个备用的默认值。
// 参数:
//
//	key: 环境变量的名称。
//	defaultValue: 如果指定的环境变量未设置，则返回的默认值。
//
// 返回值:
//
//	如果环境变量已设置并有值，则返回该值；否则返回默认值。
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetTimestampTolerance 返回时间戳容忍度，单位为秒。
// 时间戳容忍度是系统允许的时间戳误差范围，用于处理时间同步不精确的情况。
// 此函数首先尝试从环境变量中获取时间戳容忍度的值，如果获取失败或解析错误，则使用默认值300秒。
// 这种设计允许系统在不同的环境下灵活配置时间戳的容忍度，以适应不同的时间同步精度需求。
func GetTimestampTolerance() int64 {
	// 尝试从环境变量中获取时间戳容忍度的值，如果未设置，则默认为"300"
	value := getEnv("TIMESTAMP_TOLERANCE", "300")

	// 将获取到的字符串值尝试转换为int64类型
	timestampTolerance, err := strconv.ParseInt(value, 10, 64)

	// 如果转换过程中出现错误，表明环境变量的值不是有效的数字，此时返回默认的时间戳容忍度300秒
	if err != nil {
		return 300
	}

	// 如果转换成功，返回转换后的值作为时间戳容忍度
	return timestampTolerance
}

// GetGiteeSignKey 返回从环境变量中获取的 Gitee 签名密钥。
// 该函数主要用于获取与 Gitee 集成时所需的签名密钥，以便进行安全的数据交换。
// 如果环境变量中未设置 GITEE_SIGN_KEY，那么该函数将返回空字符串。
// 参数: 无
// 返回值: string 类型的签名密钥，如果未找到则为空字符串。
func GetGiteeSignKey() string {
	return getEnv("GITEE_SIGN_KEY", "")
}

// GetMaxShellExecConcurrentSize 获取最大并发执行shell命令的数量
// 该函数通过读取环境变量 "MAX_CONCURRENT" 来确定最大并发量
// 如果环境变量未设置或设置的值无效（非正整数），则默认返回 1
func GetMaxShellExecConcurrentSize() int {
	// 获取环境变量 "MAX_CONCURRENT" 的值，如果没有设置，则默认为 "1"
	value := getEnv("MAX_CONCURRENT", "1")

	// 将环境变量的值转换为整数
	maxShellExecConcurrentSize, err := strconv.Atoi(value)

	// 如果转换失败（即环境变量的值不是有效的整数），则返回默认值 1
	if err != nil {
		return 1
	}

	// 如果转换后的值小于 1（不合理的情况），也返回默认值 1
	if maxShellExecConcurrentSize < 1 {
		return 1
	}

	// 返回转换后的环境变量值作为最大并发量
	return maxShellExecConcurrentSize
}

// GetPort 返回应用程序应该监听的端口号。
// 它首先尝试从环境变量中获取 'PORT' 的值，如果未设置或为空，则默认使用 8080。
// 如果环境变量中的 'PORT' 值不是有效的整数，也将返回默认的 8080。
func GetPort() int {
	// 尝试从环境变量中获取 'PORT' 的值，如果没有设置则使用 '8080' 作为默认值。
	value := getEnv("PORT", "8080")

	// 将获取到的端口值转换为整数类型。
	port, err := strconv.Atoi(value)
	// 如果转换失败（即环境变量中的值不是有效的整数），则返回默认的端口号 8080。
	if err != nil {
		return 8080
	}

	// 如果转换成功，则返回转换后的端口号。
	return port
}

// GetLogFilePath 返回日志文件的路径。
// 该函数通过调用getEnv函数，使用环境变量"LOG_FILE_PATH"作为键来获取日志文件路径的环境变量值。
// 如果环境变量未设置，则返回默认值"logs/app.log"。
// 主要用途是提供一种灵活的方式来确定日志文件的位置，而不需要在代码中硬编码。
func GetLogFilePath() string {
	return getEnv("LOG_FILE_PATH", "logs/app.log")
}
