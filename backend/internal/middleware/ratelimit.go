package middleware

import (
	"context"
	"time"

	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/database"
	"github.com/Danche23/Evenstar-Writings/pkg/response"

	"github.com/gin-gonic/gin"
)

// RateLimit IP 限流中间件：同一 IP 在 window 内最多 limit 次（按路由 + IP 计数）
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate-limit:" + c.FullPath() + ":" + ip
		ctx := context.Background()

		count, err := database.GetRedis().Incr(ctx, key).Result()
		if err != nil {
			// Redis 异常时放行，避免限流组件故障拖垮接口
			c.Next()
			return
		}
		if count == 1 {
			_ = database.GetRedis().Expire(ctx, key, window)
		}
		if count > int64(limit) {
			response.Error(c, apperrors.CodeTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
