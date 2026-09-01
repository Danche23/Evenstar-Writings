package entity

// ArticleCategory 文章-分类中间表（联合主键，硬删除）
type ArticleCategory struct {
	ArticleID  uint `gorm:"column:article_id;primaryKey" json:"article_id"`
	CategoryID uint `gorm:"column:category_id;primaryKey" json:"category_id"`
}

// TableName 表名
func (ArticleCategory) TableName() string {
	return "article_categories"
}

// ArticleTag 文章-标签中间表（联合主键，硬删除）
type ArticleTag struct {
	ArticleID uint `gorm:"column:article_id;primaryKey" json:"article_id"`
	TagID     uint `gorm:"column:tag_id;primaryKey" json:"tag_id"`
}

// TableName 表名
func (ArticleTag) TableName() string {
	return "article_tags"
}
