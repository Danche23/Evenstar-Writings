package app

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/Danche23/Evenstar-Writings/internal/model/entity"
	"github.com/Danche23/Evenstar-Writings/pkg/logger"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// ensureAdmin 管理员初始化（仓库零敏感凭证方案）：
// 首次启动检测 users 表无 role=1 管理员时，读 EVENSTAR_ADMIN_* 环境变量创建；
// 密码未设置时生成随机强密码，仅在日志打印一次。
func (a *App) ensureAdmin() {
	var count int64
	if err := a.mysqlDB.Model(&entity.User{}).Where("role = ?", 1).Count(&count).Error; err != nil {
		logger.Warn("检测管理员账号失败", zap.Error(err))
		return
	}
	if count > 0 {
		return
	}

	username := os.Getenv("EVENSTAR_ADMIN_USERNAME")
	email := os.Getenv("EVENSTAR_ADMIN_EMAIL")
	password := os.Getenv("EVENSTAR_ADMIN_PASSWORD")

	if username == "" {
		username = "admin"
	}
	if email == "" {
		logger.Warn("未设置 EVENSTAR_ADMIN_EMAIL，跳过管理员初始化（设置后重启即可再次触发）")
		return
	}

	randomPwd := false
	if password == "" {
		password = generateRandomPassword(16)
		randomPwd = true
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("管理员密码加密失败", zap.Error(err))
		return
	}

	admin := &entity.User{
		Username: username,
		Password: string(hash),
		Email:    email,
		Role:     1,
		Status:   1,
	}
	if err := a.mysqlDB.Create(admin).Error; err != nil {
		logger.Error("创建管理员账号失败", zap.Error(err))
		return
	}

	logger.Info("=========================================")
	logger.Info("管理员账号已初始化", zap.String("username", username), zap.String("email", email))
	if randomPwd {
		// 随机密码只在日志出现这一次，请立即登录并修改密码
		logger.Info(fmt.Sprintf("初始随机密码（仅显示一次）: %s", password))
	}
	logger.Info("=========================================")
}

// generateRandomPassword 生成指定长度的随机密码（大小写字母+数字）
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}
	return string(result)
}
