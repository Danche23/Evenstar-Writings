package storage

import (
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage 本地文件存储（开发期 mock，后续替换为 OSS 实现）
type LocalStorage struct {
	baseDir string // 本地存储根目录
	baseURL string // 对外访问前缀，如 /uploads
}

// NewLocalStorage 创建本地存储
func NewLocalStorage(baseDir, baseURL string) *LocalStorage {
	return &LocalStorage{baseDir: baseDir, baseURL: strings.TrimSuffix(baseURL, "/")}
}

// Save 保存文件到本地，返回对外访问 url
func (l *LocalStorage) Save(objectKey string, data []byte) (string, error) {
	fullPath := filepath.Join(l.baseDir, filepath.FromSlash(objectKey))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		return "", err
	}
	return l.baseURL + "/" + objectKey, nil
}

// Delete 删除本地文件（url 去掉 baseURL 前缀还原 objectKey）
func (l *LocalStorage) Delete(url string) error {
	objectKey := strings.TrimPrefix(url, l.baseURL+"/")
	fullPath := filepath.Join(l.baseDir, filepath.FromSlash(objectKey))
	return os.Remove(fullPath)
}
