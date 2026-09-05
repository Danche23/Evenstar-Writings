package dto

// 通用 DTO（跨模块复用）

// CommentUser 作者/评论者简略信息（id + nickname + avatar）
type CommentUser struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// CategoryBrief 分类简略信息
type CategoryBrief struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// TagBrief 标签简略信息
type TagBrief struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// PageData 分页响应数据（泛型，复用 list 结构）
type PageData[T any] struct {
	List      []T   `json:"list"`
	Total     int64 `json:"total"`
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	TotalPage int   `json:"total_page"`
}

// NewPageData 构造分页数据
func NewPageData[T any](list []T, total int64, page, size int) *PageData[T] {
	totalPage := int(total) / size
	if int(total)%size != 0 {
		totalPage++
	}
	if list == nil {
		list = []T{}
	}
	return &PageData[T]{
		List:      list,
		Total:     total,
		Page:      page,
		PageSize:  size,
		TotalPage: totalPage,
	}
}
