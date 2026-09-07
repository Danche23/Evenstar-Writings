package dto

// Stats 后台统计（首页四卡片）
type Stats struct {
	ArticleCount int64 `json:"article_count"` // 文章总数（含草稿）
	CommentCount int64 `json:"comment_count"` // 评论总数
	UserCount    int64 `json:"user_count"`    // 用户总数
	TotalViews   int64 `json:"total_views"`   // 总浏览量
}

// PublicStats 公开站点头部统计（页脚展示）
type PublicStats struct {
	Articles int64 `json:"articles"` // 已发布文章数
	Comments int64 `json:"comments"` // 有效评论数（含一二级，不含软删）
	Views    int64 `json:"views"`    // 已发布文章总浏览量
}
