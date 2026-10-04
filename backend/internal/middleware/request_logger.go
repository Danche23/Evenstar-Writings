package middleware

import (
	"bytes"
	"io"
	"regexp"
	"time"

	"github.com/Danche23/Evenstar-Writings/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// passwordFieldRe 匹配请求体里的密码字段（password / new_password / old_password）。
// 按字段名脱敏而不是按路由白名单：以后新增接口只要带这些字段就自动生效，不会漏。
var passwordFieldRe = regexp.MustCompile(`"(password|new_password|old_password)"\s*:\s*"[^"]*"`)

// redactPasswords 把请求体里的密码值替换为 ***，其余内容保留，兼顾排查与安全
func redactPasswords(body string) string {
	return passwordFieldRe.ReplaceAllString(body, `"$1":"***"`)
}

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 读取请求体
		var body string
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if contentType != "multipart/form-data" {
				bodyBytes, _ := io.ReadAll(c.Request.Body)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				// 先脱敏再截断：否则密码字段可能正好落在截断处，留下半截明文
				body = redactPasswords(string(bodyBytes))
				if len(body) > 1000 {
					body = body[:1000] + "..."
				}
			} else {
				body = "[File Upload]"
			}
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("body", body),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		)
	}
}
