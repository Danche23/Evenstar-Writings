package repository

import (
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用户数据访问
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByEmail 按邮箱查找用户（未删除）
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 按用户名查找用户（未删除）
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 按 ID 查找用户（未删除）
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update 按 ID 更新指定字段（map，用于改密/改资料/改状态等局部更新）
func (r *UserRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

// List 分页查询用户，keyword 按用户名/昵称/邮箱模糊搜索
// FindByIDs 批量按 ID 查用户（避免逐条查库造成 N+1）
func (r *UserRepository) FindByIDs(ids []uint) ([]model.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var users []model.User
	if err := r.db.Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) List(page, size int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	q := r.db.Model(&model.User{})
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", like, like, like)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// Delete 软删用户（GORM 软删除：写入 deleted_at）
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// GetAuthor 获取站点博主（role=1 且 status=1 的第一个用户）；无则返回 nil
func (r *UserRepository) GetAuthor() (*model.User, error) {
	var user model.User
	if err := r.db.Where("role = ? AND status = ?", 1, 1).Order("id ASC").Limit(1).Find(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}
