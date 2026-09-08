package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var globalConfig *Config

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件路径
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
	}

	// 环境变量前缀
	v.SetEnvPrefix("EVENSTAR")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置
	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 从环境变量覆盖敏感配置
	if val := os.Getenv("MYSQL_PASSWORD"); val != "" {
		config.Database.MySQL.Password = val
	}
	if val := os.Getenv("REDIS_PASSWORD"); val != "" {
		config.Database.Redis.Password = val
	}
	if val := os.Getenv("JWT_SECRET"); val != "" {
		config.JWT.Secret = val
	}
	if val := os.Getenv("MAIL_PASSWORD"); val != "" {
		config.Mail.Password = val
	}
	if val := os.Getenv("OSS_ACCESS_KEY_SECRET"); val != "" {
		config.OSS.AccessKeySecret = val
	}
	if val := os.Getenv("CAPTCHA_ACCESS_KEY_SECRET"); val != "" {
		config.Captcha.AccessKeySecret = val
	}

	// 运行环境（生产用环境变量覆盖，避免为部署改动配置文件）
	if val := os.Getenv("APP_MODE"); val != "" {
		config.App.Mode = val
	}
	if val := os.Getenv("APP_PORT"); val != "" {
		if p, err := strconv.Atoi(val); err == nil && p > 0 {
			config.App.Port = p
		}
	}
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		config.Log.Level = val
	}
	if val := os.Getenv("MAIL_MOCK"); val != "" {
		config.Mail.Mock = val == "true" || val == "1"
	}
	if val := os.Getenv("MAIL_HOST"); val != "" {
		config.Mail.Host = val
	}
	if val := os.Getenv("MAIL_USERNAME"); val != "" {
		config.Mail.Username = val
	}
	if val := os.Getenv("MAIL_FROM_ADDR"); val != "" {
		config.Mail.FromAddr = val
	}
	if val := os.Getenv("OSS_ACCESS_KEY_ID"); val != "" {
		config.OSS.AccessKeyID = val
	}
	if val := os.Getenv("CAPTCHA_ACCESS_KEY_ID"); val != "" {
		config.Captcha.AccessKeyID = val
	}
	if val := os.Getenv("CAPTCHA_SCENE_ID"); val != "" {
		config.Captcha.SceneID = val
	}
	if val := os.Getenv("MYSQL_HOST"); val != "" {
		config.Database.MySQL.Host = val
	}
	if val := os.Getenv("REDIS_HOST"); val != "" {
		config.Database.Redis.Host = val
	}
	// CORS 允许来源：逗号分隔（生产填真实域名；Nginx 同域反代时可不配）
	if val := os.Getenv("CORS_ALLOW_ORIGINS"); val != "" {
		origins := make([]string, 0)
		for _, part := range strings.Split(val, ",") {
			if o := strings.TrimSpace(part); o != "" {
				origins = append(origins, o)
			}
		}
		if len(origins) > 0 {
			config.CORS.AllowOrigins = origins
		}
	}

	globalConfig = config
	return config, nil
}

// Get 获取全局配置
func Get() *Config {
	if globalConfig == nil {
		panic("配置未初始化，请先调用 Load() 加载配置")
	}
	return globalConfig
}

// MustLoad 加载配置，失败时 panic
func MustLoad(configPath string) *Config {
	config, err := Load(configPath)
	if err != nil {
		panic(err)
	}
	return config
}
