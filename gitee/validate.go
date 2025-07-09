package gitee

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"webhook/config"
)

// validateToken 验证给定的token和timestamp是否有效。
// 参数:
//
//	token: 待验证的token。
//	timestamp: 时间戳。
//
// 返回值:
//
//	int: HTTP状态码。
//	error: 错误信息，如果验证失败会返回具体的错误信息，否则返回 nil
func validateToken(token, timestamp string) (int, error) {
	// 检查token和timestamp是否为空
	if token == "" || timestamp == "" {
		return http.StatusBadRequest, errors.New("没有header头")
	}

	// 从配置中读取时间差
	timestampTolerance := config.GetTimestampTolerance()

	// 验证时间戳是否在有效期内
	currentTime := time.Now().Unix()
	timestamps, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("时间戳格式错误")
	}

	// 将毫秒转换为秒
	timestampInSeconds := timestamps / 1000

	// 检查时间戳是否在允许的范围内
	if abs(currentTime-timestampInSeconds) > int64(timestampTolerance) {
		return http.StatusUnauthorized, errors.New("时间戳不正确")
	}

	// 计算签名
	signKey := config.GetGiteeSignKey()
	secStr := fmt.Sprintf("%s\n%s", timestamp, signKey)

	mac := hmac.New(sha256.New, []byte(signKey))
	mac.Write([]byte(secStr))
	computeToken := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// 安全比较token
	if !hmac.Equal([]byte(computeToken), []byte(token)) {
		return http.StatusUnauthorized, errors.New("token不正确")
	}
	return http.StatusOK, nil
}

// abs 返回一个整数的绝对值。
// 参数：x 一个整数
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
