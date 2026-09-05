package service

import (
	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/repository"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户业务逻辑
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetProfile 获取用户信息
func (s *UserService) GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}
	resp := dto.ToUserResponse(user)
	return &resp, nil
}

// UpdateProfile 修改个人资料（nickname/avatar，只更新传入的字段）
func (s *UserService) UpdateProfile(userID uint, req dto.UpdateProfileRequest) error {
	updates := map[string]interface{}{}
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}
	if len(updates) == 0 {
		return nil
	}
	return s.userRepo.Update(userID, updates)
}

// UpdatePassword 修改密码：验旧密码 → 密码规则 → bcrypt → 更新 + token_version+1
func (s *UserService) UpdatePassword(userID uint, req dto.UpdatePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return apperrors.ErrUserNotFound
	}

	// 校验旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return apperrors.ErrOldPasswordError
	}

	// 密码规则
	if !validPassword(req.NewPassword) {
		return apperrors.New(apperrors.CodeInvalidParam, "密码需 8 位以上且包含字母和数字")
	}

	// bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.ErrInternalError
	}

	// 更新密码 + token_version+1（旧 token 全部失效）
	return s.userRepo.Update(userID, map[string]interface{}{
		"password":      string(hash),
		"token_version": user.TokenVersion + 1,
	})
}

// List 后台分页查询用户
func (s *UserService) List(page, size int, keyword string) (*dto.PageData[dto.UserResponse], error) {
	users, total, err := s.userRepo.List(page, size, keyword)
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	list := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		list = append(list, dto.ToUserResponse(&users[i]))
	}
	return dto.NewPageData(list, total, page, size), nil
}

// Delete 后台删除用户（软删 + token_version+1 强制下线；不能删自己）
func (s *UserService) Delete(operatorID, targetID uint) error {
	if operatorID == targetID {
		return apperrors.ErrCannotOperateSelf
	}
	user, err := s.userRepo.FindByID(targetID)
	if err != nil {
		return apperrors.ErrUserNotFound
	}
	// 先 token_version+1 强制下线，再软删
	if err := s.userRepo.Update(targetID, map[string]interface{}{"token_version": user.TokenVersion + 1}); err != nil {
		return apperrors.ErrInternalError
	}
	return s.userRepo.Delete(targetID)
}

// UpdateStatus 后台禁用/解禁用户（禁用时 token_version+1 强制下线；不能操作自己）
func (s *UserService) UpdateStatus(operatorID, targetID uint, status int8) error {
	if operatorID == targetID {
		return apperrors.ErrCannotOperateSelf
	}
	user, err := s.userRepo.FindByID(targetID)
	if err != nil {
		return apperrors.ErrUserNotFound
	}
	updates := map[string]interface{}{"status": status}
	if status == 2 {
		updates["token_version"] = user.TokenVersion + 1
	}
	return s.userRepo.Update(targetID, updates)
}
