package dto

// Stats 后台统计（首页四卡片）
type Stats struct {
	ArticleCount int64 `json:"article_count"` // 文章总数（含草稿）
	CommentCount int64 `json:"comment_count"` // 评论总数
	UserCount    int64 `json:"user_count"`    // 用户总数
	TotalViews   int64 `json:"total_views"`   // 总浏览量
}
