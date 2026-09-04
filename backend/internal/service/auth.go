package service

// 认证业务逻辑层。
import (
	"errors"

	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/repository"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 认证业务逻辑
type AuthService struct {
	userRepo *repository.UserRepository
}

// 创建认证服务
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Login 登录：查用户 → 验密码 → 验状态 → 发 JWT
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. 按邮箱查用户
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, apperrors.ErrInternalError
	}

	// 2. 校验密码(bcrypt比对)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	// 3. 校验账号状态(1=正常 2=禁用)
	if user.Status != 1 {
		return nil, apperrors.ErrUserDisabled
	}

	// 4. 生成jwt
	token, err := jwt.GenerateToken(user.ID, user.Username, user.TokenVersion)
	if err != nil {
		return nil, apperrors.ErrInternalError
	}

	return &dto.LoginResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil
}
