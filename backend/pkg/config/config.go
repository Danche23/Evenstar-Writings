package config

// Config 应用配置结构体
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Mail     MailConfig     `mapstructure:"mail"`
	OSS      OSSConfig      `mapstructure:"oss"`
	Captcha  CaptchaConfig  `mapstructure:"captcha"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Mode    string `mapstructure:"mode"` // debug, release, test
	Port    int    `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
}

// MySQLConfig MySQL 配置
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"` // 有效期（小时），设计定稿 168 = 7 天
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Filename   string `mapstructure:"filename"`    // 日志文件路径
	MaxSize    int    `mapstructure:"max_size"`    // 单个日志文件最大大小(MB)
	MaxBackups int    `mapstructure:"max_backups"` // 保留的旧日志文件数量
	MaxAge     int    `mapstructure:"max_age"`     // 保留旧日志文件的最大天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩
}

// CORSConfig CORS 配置
type CORSConfig struct {
	Enabled          bool     `mapstructure:"enabled"`
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// MailConfig 邮件配置（阿里云邮件推送）
type MailConfig struct {
	Mock     bool   `mapstructure:"mock"`      // true 时不真发邮件，验证码打印到日志（开发期）
	Host     string `mapstructure:"host"`      // SMTP 地址
	Port     int    `mapstructure:"port"`      // SMTP 端口
	Username string `mapstructure:"username"`  // SMTP 账号
	Password string `mapstructure:"password"`  // SMTP 授权码（建议走环境变量 MAIL_PASSWORD）
	FromAddr string `mapstructure:"from_addr"` // 发信地址
	FromName string `mapstructure:"from_name"` // 发信人名称
}

// OSSConfig 阿里云 OSS 配置
type OSSConfig struct {
	Bucket          string `mapstructure:"bucket"`
	Region          string `mapstructure:"region"` // Bucket 所在地域，如 cn-hangzhou；endpoint 留空时自动拼 oss-{region}.aliyuncs.com
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"` // 建议走环境变量 OSS_ACCESS_KEY_SECRET
}

// CaptchaConfig 阿里云验证码 2.0 配置
type CaptchaConfig struct {
	SceneID         string `mapstructure:"scene_id"`        // 验证场景 ID
	IdentityPrefix  string `mapstructure:"identity_prefix"` // 身份标（前端初始化 AliyunCaptcha 时使用）
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"` // 建议走环境变量 CAPTCHA_ACCESS_KEY_SECRET
	Endpoint        string `mapstructure:"endpoint"`          // captcha.cn-shanghai.aliyuncs.com
}
