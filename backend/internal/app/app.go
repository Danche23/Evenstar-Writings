package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Danche23/Evenstar-Writings/internal/api"
	"github.com/Danche23/Evenstar-Writings/internal/api/article"
	"github.com/Danche23/Evenstar-Writings/internal/api/auth"
	"github.com/Danche23/Evenstar-Writings/internal/api/category"
	"github.com/Danche23/Evenstar-Writings/internal/api/comment"
	"github.com/Danche23/Evenstar-Writings/internal/api/stats"
	"github.com/Danche23/Evenstar-Writings/internal/api/tag"
	"github.com/Danche23/Evenstar-Writings/internal/api/upload"
	"github.com/Danche23/Evenstar-Writings/internal/api/user"
	"github.com/Danche23/Evenstar-Writings/internal/repository"
	"github.com/Danche23/Evenstar-Writings/internal/service"
	"github.com/Danche23/Evenstar-Writings/pkg/config"
	"github.com/Danche23/Evenstar-Writings/pkg/database"
	"github.com/Danche23/Evenstar-Writings/pkg/logger"
	"github.com/Danche23/Evenstar-Writings/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App 应用结构体
type App struct {
	cfg            *config.Config
	mysqlDB        *gorm.DB
	redis          *redis.Client
	router         *api.Router
	server         *http.Server
	articleService *service.ArticleService
}

// NewApp 创建应用实例
func NewApp() *App {
	return &App{}
}

// Initialize 初始化应用
func (a *App) Initialize(configPath string) error {
	// 1. 加载配置
	if err := a.initConfig(configPath); err != nil {
		return err
	}

	// 1.5 关键配置校验：JWT 密钥非空，否则启动即失败
	if a.cfg.JWT.Secret == "" {
		return fmt.Errorf("JWT 密钥未配置：请在 config.yaml 或环境变量中设置 jwt.secret")
	}

	// 2. 初始化日志
	if err := a.initLogger(); err != nil {
		return err
	}

	// 3. 初始化数据库
	if err := a.initDatabase(); err != nil {
		return err
	}

	// 4. 初始化依赖
	a.initDependencies()

	// 4.5 启动定时任务（浏览量回写）
	a.startCron()

	// 5. 管理员初始化（首次启动且无管理员时创建）
	a.ensureAdmin()

	// 6. 初始化路由
	a.initRouter()

	// 7. 初始化服务器
	a.initServer()

	return nil
}

// initConfig 加载配置
func (a *App) initConfig(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	a.cfg = cfg
	return nil
}

// initLogger 初始化日志
func (a *App) initLogger() error {
	if err := logger.Init(&a.cfg.Log); err != nil {
		return fmt.Errorf("日志初始化失败: %w", err)
	}

	// 打印启动横幅
	logger.Info("=========================================")
	logger.Info(fmt.Sprintf("欢迎使用 %s", a.cfg.App.Name))
	logger.Info(fmt.Sprintf("版本: %s", a.cfg.App.Version))
	logger.Info(fmt.Sprintf("模式: %s", a.cfg.App.Mode))
	logger.Info("配置加载成功")
	logger.Info("=========================================")

	return nil
}

// initDatabase 初始化数据库
func (a *App) initDatabase() error {
	// 初始化 MySQL
	mysqlDB, err := database.InitMySQL(&a.cfg.Database.MySQL)
	if err != nil {
		return fmt.Errorf("MySQL 初始化失败: %w", err)
	}
	a.mysqlDB = mysqlDB

	// 表结构以 scripts/init.sql 为准（含外键/索引），不在代码里 AutoMigrate

	// 初始化 Redis（可选，失败不影响核心功能）
	rs, err := database.InitRedis(&a.cfg.Database.Redis)
	if err != nil {
		logger.Warn("Redis 初始化失败，将不影响核心功能", zap.Error(err))
	}
	a.redis = rs

	return nil
}

// initDependencies 初始化依赖注入
func (a *App) initDependencies() error {
	// ========== 创建 Repository ==========
	userRepo := repository.NewUserRepository(a.mysqlDB)
	articleRepo := repository.NewArticleRepository(a.mysqlDB)
	categoryRepo := repository.NewCategoryRepository(a.mysqlDB)
	tagRepo := repository.NewTagRepository(a.mysqlDB)
	commentRepo := repository.NewCommentRepository(a.mysqlDB)
	uploadRepo := repository.NewUploadRepository(a.mysqlDB)

	// ========== 存储实现（OSS 配置填齐则启用正式存储，否则本地 mock） ==========
	uploadStorage, err := a.newUploadStorage()
	if err != nil {
		return err
	}

	// ========== 创建 Service ==========
	mailService := service.NewMailService(&a.cfg.Mail)
	captchaService := service.NewCaptchaService(&a.cfg.Captcha)
	authService := service.NewAuthService(userRepo, a.redis, mailService, captchaService)
	userService := service.NewUserService(userRepo)
	articleService := service.NewArticleService(articleRepo, userRepo, categoryRepo, tagRepo, a.redis)
	a.articleService = articleService
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	commentService := service.NewCommentService(commentRepo, userRepo, articleRepo, a.redis)
	uploadService := service.NewUploadService(uploadRepo, articleRepo, uploadStorage, a.redis)
	statsService := service.NewStatsService(a.mysqlDB)

	// ========== 创建 Handler ==========
	authHandler := auth.NewAuthHandler(authService)
	userHandler := user.NewUserHandler(userService)
	articleHandler := article.NewArticleHandler(articleService)
	categoryHandler := category.NewCategoryHandler(categoryService)
	tagHandler := tag.NewTagHandler(tagService)
	commentHandler := comment.NewCommentHandler(commentService)
	uploadHandler := upload.NewUploadHandler(uploadService)
	statsHandler := stats.NewStatsHandler(statsService)

	// ========== 创建 Router ==========
	a.router = api.NewRouter(authHandler, userHandler, articleHandler, categoryHandler, tagHandler, commentHandler, uploadHandler, statsHandler)
	return nil
}

// newUploadStorage 根据 oss 配置选择存储实现：
//   - bucket + endpoint(或 region) + access_key_id + access_key_secret 四项填齐 → 阿里云 OSS（初始化失败则启动失败，不静默回退）
//   - 未填齐 → 本地 mock（开发期），仅提示缺哪些字段，不打印任何密钥
func (a *App) newUploadStorage() (storage.Storage, error) {
	ossCfg := &a.cfg.OSS
	configured := ossCfg.Bucket != "" &&
		(ossCfg.Endpoint != "" || ossCfg.Region != "") &&
		ossCfg.AccessKeyID != "" &&
		ossCfg.AccessKeySecret != ""
	if !configured {
		logger.Info("OSS 未启用（bucket / region|endpoint / access_key 未填齐），文件存储使用本地 mock：./storage")
		return storage.NewLocalStorage("./storage", "/uploads"), nil
	}
	st, err := storage.NewOSSStorage(ossCfg)
	if err != nil {
		return nil, fmt.Errorf("OSS 初始化失败: %w", err)
	}
	logger.Info("已启用阿里云 OSS 文件存储")
	return st, nil
}

// startCron 启动定时任务（浏览量回写，每 5 分钟）
func (a *App) startCron() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if a.articleService != nil {
				if err := a.articleService.SyncViews(); err != nil {
					logger.Error("浏览量回写失败", zap.Error(err))
				}
			}
		}
	}()
}

// initRouter 初始化路由
func (a *App) initRouter() {
	// 设置 Gin 模式
	gin.SetMode(a.cfg.App.Mode)
}

// initServer 初始化 HTTP 服务器
func (a *App) initServer() {
	engine := gin.New()

	// 注册路由
	a.router.Setup(engine)

	// 创建 HTTP 服务器
	a.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", a.cfg.App.Port),
		Handler:        engine,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
}

// Run 运行应用
func (a *App) Run() {
	// 启动 HTTP 服务器
	go func() {
		logger.Info("HTTP 服务器启动",
			zap.String("addr", a.server.Addr),
			zap.String("mode", a.cfg.App.Mode),
		)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP 服务器启动失败", zap.Error(err))
		}
	}()

	// 优雅关闭
	a.gracefulShutdown()
}

// gracefulShutdown 优雅关闭
func (a *App) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭 HTTP 服务器
	if err := a.server.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", zap.Error(err))
	}

	// 关闭路由连接
	if a.router != nil {
		if err := a.router.Close(); err != nil {
			logger.Error("关闭路由连接失败", zap.Error(err))
		}
	}

	// 关闭数据库连接
	_ = database.CloseMySQL()
	_ = database.CloseRedis()

	// 同步日志
	_ = logger.Sync()

	logger.Info("服务器已关闭")
	logger.Info("=========================================")
}
