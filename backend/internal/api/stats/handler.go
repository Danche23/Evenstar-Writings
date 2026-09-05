package stats

import (
	"github.com/Danche23/Evenstar-Writings/internal/service"
	"github.com/Danche23/Evenstar-Writings/pkg/response"

	"github.com/gin-gonic/gin"
)

// StatsHandler 统计模块处理器
type StatsHandler struct {
	statsService *service.StatsService
}

// NewStatsHandler 创建统计处理器
func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// Stats 后台统计
func (h *StatsHandler) Stats(c *gin.Context) {
	resp, err := h.statsService.Stats()
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, resp)
}
