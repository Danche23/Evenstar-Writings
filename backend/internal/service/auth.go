package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Danche23/Evenstar-Writings/internal/dto"
	"github.com/Danche23/Evenstar-Writings/internal/model"
	"github.com/Danche23/Evenstar-Writings/internal/repository"
	apperrors "github.com/Danche23/Evenstar-Writings/pkg/errors"
	"github.com/Danche23/Evenstar-Writings/pkg/jwt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证业务逻辑
type AuthService struct {
	userRepo *repository.UserRepository
	redis    *redis.Client
	mail     *MailService
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *repository.UserRepository, rdb *redis.Client, mail *MailService) *AuthService {
	return &AuthService{userRepo: userRepo, redis: rdb, mail: mail}
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

	// 2. 校验密码(bcrypt 比对)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	// 3. 校验账号状态(1=正常 2=禁用)
	if user.Status != 1 {
		return nil, apperrors.ErrUserDisabled
	}

	// 4. 生成 jwt
	token, err := jwt.GenerateToken(user.ID, user.Username, user.TokenVersion)
	if err != nil {
		return nil, apperrors.ErrInternalError
	}

	return &dto.LoginResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil
}

// SendCode 发送邮箱验证码（register：邮箱未注册；reset：邮箱已注册）
func (s *AuthService) SendCode(req dto.SendCodeRequest) error {
	ctx := context.Background()

	// 1. 按 type 校验邮箱状态
	_, err := s.userRepo.FindByEmail(req.Email)
	if req.Type == "register" {
		// 注册：邮箱必须未注册
		if err == nil {
			return apperrors.ErrUserAlreadyExists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrInternalError
		}
	} else {
		// reset：邮箱必须已注册
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrUserNotFound
		}
		if err != nil {
			return apperrors.ErrInternalError
		}
	}

	// 2. 频率限制：同邮箱 60 秒一次
	limitKey := fmt.Sprintf("send-code:limit:%s", req.Email)
	ok, err := s.redis.SetNX(ctx, limitKey, 1, 60*time.Second).Result()
	if err != nil {
		return apperrors.ErrInternalError
	}
	if !ok {
		return apperrors.New(apperrors.CodeTooManyRequests, "发送太频繁，请稍后再试")
	}

	// 3. 频率限制：每日 10 次
	dailyKey := fmt.Sprintf("send-code:daily:%s:%s", req.Email, time.Now().Format("2006-01-02"))
	count, err := s.redis.Incr(ctx, dailyKey).Result()
	if err != nil {
		return apperrors.ErrInternalError
	}
	if count == 1 {
		_ = s.redis.Expire(ctx, dailyKey, 24*time.Hour)
	}
	if count > 10 {
		return apperrors.New(apperrors.CodeTooManyRequests, "今日发送次数已达上限")
	}

	// 4. 生成验证码并存储（5 分钟）
	code := generateCode()
	codeKey := fmt.Sprintf("register:code:%s", req.Email)
	if err := s.redis.Set(ctx, codeKey, code, 5*time.Minute).Err(); err != nil {
		return apperrors.ErrInternalError
	}

	// 5. 发邮件（mock 时打日志）
	subject := "【Evenstar Writings】验证码"
	body := fmt.Sprintf("您的验证码是：%s，5 分钟内有效。", code)
	if err := s.mail.Send(req.Email, subject, body); err != nil {
		return apperrors.ErrInternalError
	}

	return nil
}

// Register 注册：校验验证码 → 唯一性 → 密码规则 → bcrypt → 创建 → 自动登录
func (s *AuthService) Register(req dto.RegisterRequest) (*dto.LoginResponse, error) {
	ctx := context.Background()

	// 1. 校验验证码
	codeKey := fmt.Sprintf("register:code:%s", req.Email)
	if err := s.verifyCode(ctx, codeKey, req.Code); err != nil {
		return nil, err
	}

	// 2. 用户名唯一
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, apperrors.ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrInternalError
	}

	// 3. 邮箱唯一
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, apperrors.ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrInternalError
	}

	// 4. 密码规则（8 位以上含字母和数字）
	if !validPassword(req.Password) {
		return nil, apperrors.New(apperrors.CodeInvalidParam, "密码需 8 位以上且包含字母和数字")
	}

	// 5. bcrypt 哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.ErrInternalError
	}

	// 6. 创建用户（role 恒 2，昵称空则用用户名兜底）
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}
	user := &model.User{
		Username: req.Username,
		Password: string(hash),
		Nickname: nickname,
		Email:    req.Email,
		Role:     2,
		Status:   1,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, apperrors.ErrInternalError
	}

	// 7. 删除已用验证码
	_ = s.redis.Del(ctx, codeKey)

	// 8. 生成 token 自动登录
	token, err := jwt.GenerateToken(user.ID, user.Username, user.TokenVersion)
	if err != nil {
		return nil, apperrors.ErrInternalError
	}

	return &dto.LoginResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil
}

// ResetPassword 找回密码：校验验证码 → 查用户 → 密码规则 → bcrypt → 更新 + token_version+1
func (s *AuthService) ResetPassword(req dto.ResetPasswordRequest) error {
	ctx := context.Background()

	// 1. 校验验证码
	codeKey := fmt.Sprintf("register:code:%s", req.Email)
	if err := s.verifyCode(ctx, codeKey, req.Code); err != nil {
		return err
	}

	// 2. 查用户
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrUserNotFound
		}
		return apperrors.ErrInternalError
	}

	// 3. 密码规则
	if !validPassword(req.NewPassword) {
		return apperrors.New(apperrors.CodeInvalidParam, "密码需 8 位以上且包含字母和数字")
	}

	// 4. bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.ErrInternalError
	}

	// 5. 更新密码 + token_version+1（旧 token 全部失效）
	if err := s.userRepo.Update(user.ID, map[string]interface{}{
		"password":      string(hash),
		"token_version": user.TokenVersion + 1,
	}); err != nil {
		return apperrors.ErrInternalError
	}

	// 6. 删除已用验证码
	_ = s.redis.Del(ctx, codeKey)

	return nil
}

// verifyCode 校验验证码（从 Redis 取，比对，不匹配返回业务错误）
func (s *AuthService) verifyCode(ctx context.Context, key, code string) error {
	got, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return apperrors.ErrVerifyCodeError
	}
	if err != nil {
		return apperrors.ErrInternalError
	}
	if got != code {
		return apperrors.ErrVerifyCodeError
	}
	return nil
}

// generateCode 生成 6 位数字验证码
func generateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// validPassword 校验密码规则：8 位以上且包含字母和数字
func validPassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}
	var hasLetter, hasDigit bool
	for _, c := range pw {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}
