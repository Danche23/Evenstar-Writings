package repository

import (
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"gorm.io/gorm"
)

// ArticleRepository 文章数据访问
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓库
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// ListPublished 前台文章列表（仅已发布 status=2，支持分类/标签/关键词筛选）
func (r *ArticleRepository) ListPublished(page, size int, categoryID, tagID uint, keyword string) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	q := r.db.Model(&model.Article{}).Where("status = ?", 2)
	if categoryID > 0 {
		q = q.Where("id IN (SELECT article_id FROM article_categories WHERE category_id = ?)", categoryID)
	}
	if tagID > 0 {
		q = q.Where("id IN (SELECT article_id FROM article_tags WHERE tag_id = ?)", tagID)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("title LIKE ? OR summary LIKE ?", like, like)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("published_at DESC").Offset((page - 1) * size).Limit(size).Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

// List 后台文章列表（含草稿，status 可选、keyword 标题搜索）
func (r *ArticleRepository) List(page, size int, status int8, keyword string) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	q := r.db.Model(&model.Article{})
	if status > 0 {
		q = q.Where("status = ?", status)
	}
	if keyword != "" {
		q = q.Where("title LIKE ?", "%"+keyword+"%")
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

// FindByID 按 ID 查文章（未删除）
func (r *ArticleRepository) FindByID(id uint) (*model.Article, error) {
	var article model.Article
	if err := r.db.First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// FindByIDs 批量查文章（热门榜用，保留传入顺序由调用方处理）
func (r *ArticleRepository) FindByIDs(ids []uint) ([]model.Article, error) {
	if len(ids) == 0 {
		return []model.Article{}, nil
	}
	var articles []model.Article
	if err := r.db.Where("id IN ?", ids).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// Create 创建文章
func (r *ArticleRepository) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

// Update 按 ID 更新指定字段
func (r *ArticleRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 软删文章
func (r *ArticleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Article{}, id).Error
}

// IncrViews 浏览量累加（cron 回写用，原子自增）
func (r *ArticleRepository) IncrViews(id uint, delta int64) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", delta)).Error
}

// FindByContentURL 查内容引用了指定 url 的文章（含软删，上传删除引用检查用）
func (r *ArticleRepository) FindByContentURL(url string) ([]model.Article, error) {
	var articles []model.Article
	if err := r.db.Unscoped().Where("content LIKE ?", "%"+url+"%").Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// GetCategoryIDs 查文章已关联的分类 id
func (r *ArticleRepository) GetCategoryIDs(articleID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&model.ArticleCategory{}).Where("article_id = ?", articleID).Pluck("category_id", &ids).Error
	return ids, err
}

// GetTagIDs 查文章已关联的标签 id
func (r *ArticleRepository) GetTagIDs(articleID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&model.ArticleTag{}).Where("article_id = ?", articleID).Pluck("tag_id", &ids).Error
	return ids, err
}

// SetCategories 全量替换文章分类（先删后插，事务）
func (r *ArticleRepository) SetCategories(articleID uint, categoryIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id = ?", articleID).Delete(&model.ArticleCategory{}).Error; err != nil {
			return err
		}
		for _, cid := range categoryIDs {
			if err := tx.Create(&model.ArticleCategory{ArticleID: articleID, CategoryID: cid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// SetTags 全量替换文章标签（先删后插，事务）
func (r *ArticleRepository) SetTags(articleID uint, tagIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id = ?", articleID).Delete(&model.ArticleTag{}).Error; err != nil {
			return err
		}
		for _, tid := range tagIDs {
			if err := tx.Create(&model.ArticleTag{ArticleID: articleID, TagID: tid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
