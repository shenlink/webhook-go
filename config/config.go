package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv 加载 .env 文件中的配置
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("加载 .env 文件失败: %v", err)
	}
}

// getEnv 获取指定键的值，如果不存在则返回默认值
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetTimestampTolerance() int64 {
	value := getEnv("TIMESTAMP_TOLERANCE", "300")
	timestampTolerance, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 300
	}
	return timestampTolerance
}

func GetGiteeSignKey() string {
	return getEnv("GITEE_SIGN_KEY", "")
}

func GetMaxShellExecConcurrentSize() int {
	value := getEnv("MAX_CONCURRENT", "1")
	maxShellExecConcurrentSize, err := strconv.Atoi(value)
	if err != nil {
		return 1
	}
	if maxShellExecConcurrentSize < 1 {
		return 1
	}
	return maxShellExecConcurrentSize
}

func GetPort() int {
	value := getEnv("PORT", "8080")
	port, err := strconv.Atoi(value)
	if err != nil {
		return 8080
	}
	return port
}

func GetLogFilePath() string {
	return getEnv("LOG_FILE_PATH", "logs/app.log")
}
