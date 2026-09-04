package storage

import (
	"bytes"
	"fmt"
	"mime"
	"path/filepath"
	"strings"

	"github.com/Danche23/Evenstar-Writings/pkg/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSStorage 阿里云 OSS 存储（正式实现）。
// 使用说明：config.yaml 中 oss 四项（bucket + endpoint/region + access_key_id + access_key_secret）填齐后，
// 由 app 层创建本实现；未填齐时 app 层回退 LocalStorage（开发期 mock）。
type OSSStorage struct {
	bucket    *oss.Bucket
	publicURL string // 对外访问前缀，如 https://{bucket}.{endpoint}
}

// NewOSSStorage 创建 OSS 存储（配置不完整或初始化失败返回 error，由调用方决定回退或终止）
func NewOSSStorage(cfg *config.OSSConfig) (Storage, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("OSS 配置缺失：bucket 为空")
	}
	if cfg.Endpoint == "" && cfg.Region == "" {
		return nil, fmt.Errorf("OSS 配置缺失：endpoint 与 region 均为空")
	}
	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("OSS 配置缺失：access_key_id / access_key_secret 为空")
	}

	endpoint := normalizeEndpoint(cfg.Endpoint, cfg.Region)
	client, err := oss.New(endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建 OSS 客户端失败: %w", err)
	}
	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("获取 OSS Bucket 失败: %w", err)
	}

	publicHost := strings.TrimPrefix(endpoint, "https://")
	publicHost = strings.TrimPrefix(publicHost, "http://")
	publicHost = strings.TrimSuffix(publicHost, "/")

	return &OSSStorage{
		bucket:    bucket,
		publicURL: "https://" + cfg.Bucket + "." + publicHost,
	}, nil
}

// Save 上传对象到 OSS（content-type 按扩展名自动推断），返回公网可访问 url
func (o *OSSStorage) Save(objectKey string, data []byte) (string, error) {
	opts := make([]oss.Option, 0, 1)
	if ct := mime.TypeByExtension(filepath.Ext(objectKey)); ct != "" {
		opts = append(opts, oss.ContentType(ct))
	}
	if err := o.bucket.PutObject(objectKey, bytes.NewReader(data), opts...); err != nil {
		return "", err
	}
	return o.publicURL + "/" + objectKey, nil
}

// Delete 删除对象（url 去掉 publicURL 前缀还原 objectKey）
func (o *OSSStorage) Delete(url string) error {
	objectKey := strings.TrimPrefix(url, o.publicURL+"/")
	if objectKey == url {
		return fmt.Errorf("非法文件地址: %s", url)
	}
	return o.bucket.DeleteObject(objectKey)
}

// normalizeEndpoint 规范化 endpoint：region 兜底 + 补全 https scheme
func normalizeEndpoint(endpoint, region string) string {
	e := strings.TrimSpace(endpoint)
	if e == "" {
		region = strings.TrimSpace(region)
		if strings.HasPrefix(region, "oss-") {
			e = region + ".aliyuncs.com"
		} else {
			e = "oss-" + region + ".aliyuncs.com"
		}
	}
	if !strings.HasPrefix(e, "http://") && !strings.HasPrefix(e, "https://") {
		e = "https://" + e
	}
	return strings.TrimSuffix(e, "/")
}
