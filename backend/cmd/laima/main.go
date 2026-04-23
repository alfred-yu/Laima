
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	repoapi "laima/internal/repo/api"
	repoapp "laima/internal/repo/app"
	userapi "laima/internal/user/api"
	prapi "laima/internal/pr/api"
	prapp "laima/internal/pr/app"
	aiapi "laima/internal/ai/api"
	aiapp "laima/internal/ai/app"
	cicdapi "laima/internal/cicd/api"
	issueapi "laima/internal/issue/api"
	auditapi "laima/internal/audit/api"
	"laima/internal/git"
	"laima/internal/middleware"
	"laima/internal/user/app"
	"laima/internal/user/domain"
	repodomain "laima/internal/repo/domain"
	prdomain "laima/internal/pr/domain"
	issuedomain "laima/internal/issue/domain"
	cicddomain "laima/internal/cicd/domain"
	aidomain "laima/internal/ai/domain"
	auditdomain "laima/internal/audit/domain"
	auditapp "laima/internal/audit/app"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/meilisearch/meilisearch-go"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Laima API
// @version 1.0
// @description Laima 代码托管平台 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http https


func main() {
	// 初始化数据库连接 - 暂时设为可选
	var db *gorm.DB
	var err error
	db, err = initDatabase()
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v, running without database", err)
		// 继续运行，即使没有数据库
	} else {
		// 自动迁移数据库表
		err = autoMigrate(db)
		if err != nil {
			log.Printf("Warning: Failed to migrate database: %v, running without database migrations", err)
		}
	}

	// 初始化 Redis 连接 - 可选，如果失败则使用内存存储
	redisClient, err := initRedis()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v, using memory storage instead", err)
		redisClient = nil
	}

	// 初始化 MinIO 客户端 - 可选，如果失败则使用本地存储
	minioClient, err := initMinIO()
	if err != nil {
		log.Printf("Warning: Failed to connect to MinIO: %v, using local storage instead", err)
		minioClient = nil
	}

	// 初始化 Meilisearch 客户端 - 可选，如果失败则搜索功能受限
	var meiliClient meilisearch.ServiceManager
	meiliClient, err = initMeilisearch()
	if err != nil {
		log.Printf("Warning: Failed to connect to Meilisearch: %v, search functionality will be limited", err)
		meiliClient = nil
	}

	// 初始化用户服务
	var userService app.UserService
	if db != nil {
		userService = app.NewUserService(db)
	} else {
		log.Printf("Warning: Running without user service due to missing database")
	}

	// 初始化 Git 服务
	gitSvc := initGitService()

	// 初始化仓库服务
	if db != nil {
		repoService := repoapp.NewRepoService(db, gitSvc, meiliClient)
		_ = repoService // 暂时使用变量，避免未使用的警告
	} else {
		log.Printf("Warning: Running without repo service due to missing database")
	}

	// 初始化并启动 SSH 服务器
	sshServer := initSSHServer(gitSvc, userService)
	go func() {
		if err := sshServer.Start(context.Background()); err != nil {
			log.Printf("SSH 服务器启动失败: %v", err)
		}
	}()

	// 设置 Gin 模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 Gin 引擎
	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 注册审计中间件（在认证之后，路由之前）
	if db != nil {
		auditSvc := auditapp.NewAuditService(db)
		auditMiddleware := middleware.NewAuditMiddleware(auditSvc)
		r.Use(auditMiddleware.Handler())
	} else {
		log.Printf("Warning: Running without audit service due to missing database")
	}

	// 初始化 PR 服务
	var prService prapp.PRService
	if db != nil {
		prService = prapp.NewPRService(db, gitSvc)
	}

	// 初始化 AI 服务
	var aiService aiapp.AIService
	if db != nil {
		aiService = aiapp.NewAIService(db, gitSvc, prService)
	}

	// 注册 API 路由
	if db != nil {
		repoAPI := repoapi.NewRepoAPI(db, redisClient, minioClient, meiliClient, gitSvc)
		repoAPI.RegisterRoutes(r)

		userAPI := userapi.NewUserAPI(db, redisClient)
		userAPI.RegisterRoutes(r)

		prAPI := prapi.NewPRAPI(db, gitSvc)
		prAPI.RegisterRoutes(r)

		aiAPI := aiapi.NewAIApi(db, aiService)
		aiAPI.RegisterRoutes(r)

		cicdAPI := cicdapi.NewCICDApi(db)
		cicdAPI.RegisterRoutes(r)

		issueAPI := issueapi.NewIssueApi(db)
		issueAPI.RegisterRoutes(r)

		// 注册审计路由
		auditAPI := auditapi.NewAuditAPI(db)
		auditAPI.RegisterRoutes(r)
	} else {
		log.Printf("Warning: Running without API routes due to missing database")
	}

	// 注册 Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Laima is running",
		})
	})

	// 启动服务器
	port := os.Getenv("LAIMA_HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Organization{},
		&domain.OrganizationMember{},
		&domain.RepositoryMember{},
		&repodomain.Repository{},
		&prdomain.PullRequest{},
		&prdomain.Review{},
		&prdomain.ReviewComment{},
		&issuedomain.Issue{},
		&issuedomain.IssueComment{},
		&issuedomain.Milestone{},
		&cicddomain.Pipeline{},
		&cicddomain.Job{},
		&aidomain.AIReview{},
		&aidomain.AIReviewIssue{},
		&auditdomain.AuditLog{},
	)
}

// initGitService 初始化 Git 服务
func initGitService() *git.Service {
	// 获取仓库存储路径
	repoPath := os.Getenv("LAIMA_REPO_PATH")
	if repoPath == "" {
		// 默认路径为当前工作目录下的 repos 文件夹
		cwd, _ := os.Getwd()
		repoPath = filepath.Join(cwd, "repos")
	}

	// 确保目录存在
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		log.Fatalf("Failed to create repo directory: %v", err)
	}

	return git.NewService(repoPath)
}

// initSSHServer 初始化 SSH 服务器
func initSSHServer(gitSvc *git.Service, userService app.UserService) *git.SSHServer {
	// 获取 SSH 服务器配置
	sshAddr := os.Getenv("LAIMA_SSH_PORT")
	if sshAddr == "" {
		sshAddr = "2222"
	}

	hostKeyPath := os.Getenv("LAIMA_SSH_HOST_KEY")
	if hostKeyPath == "" {
		// 默认路径为当前工作目录下的 ssh 文件夹
		cwd, _ := os.Getwd()
		hostKeyPath = filepath.Join(cwd, "ssh", "host_key")
	}

	// 确保 SSH 目录存在
	if err := os.MkdirAll(filepath.Dir(hostKeyPath), 0755); err != nil {
		log.Fatalf("Failed to create SSH directory: %v", err)
	}

	// 创建并返回 SSH 服务器
	return git.NewSSHServer(":"+sshAddr, hostKeyPath, gitSvc.GetRepoBasePath(), gitSvc, userService)
}

// 初始化数据库连接
func initDatabase() (*gorm.DB, error) {
	dbHost := os.Getenv("LAIMA_DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("LAIMA_DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbName := os.Getenv("LAIMA_DB_NAME")
	if dbName == "" {
		dbName = "laima"
	}

	dbUser := os.Getenv("LAIMA_DB_USER")
	if dbUser == "" {
		dbUser = "laima"
	}

	dbPassword := os.Getenv("LAIMA_DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "laima123"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// 初始化 Redis 连接
func initRedis() (*redis.Client, error) {
	redisURL := os.Getenv("LAIMA_REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	// 测试连接
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// 初始化 MinIO 客户端
func initMinIO() (*minio.Client, error) {
	minioEndpoint := os.Getenv("LAIMA_MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "localhost:9000"
	}

	minioAccessKey := os.Getenv("LAIMA_MINIO_ACCESS_KEY")
	if minioAccessKey == "" {
		minioAccessKey = "laima"
	}

	minioSecretKey := os.Getenv("LAIMA_MINIO_SECRET_KEY")
	if minioSecretKey == "" {
		minioSecretKey = "laima123"
	}

	useSSL := os.Getenv("LAIMA_MINIO_USE_SSL") == "true"

	client, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// 测试连接
	ctx := context.Background()
	_, err = client.BucketExists(ctx, "laima")
	if err != nil {
		// 如果 bucket 不存在，创建它
		err = client.MakeBucket(ctx, "laima", minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

// 初始化 Meilisearch 客户端
func initMeilisearch() (meilisearch.ServiceManager, error) {
	meiliURL := os.Getenv("LAIMA_MEILISEARCH_URL")
	if meiliURL == "" {
		meiliURL = "http://localhost:7700"
	}

	var options []meilisearch.Option

	meiliAPIKey := os.Getenv("LAIMA_MEILISEARCH_API_KEY")
	if meiliAPIKey != "" {
		options = append(options, meilisearch.WithAPIKey(meiliAPIKey))
	}

	client, err := meilisearch.Connect(meiliURL, options...)
	if err != nil {
		return nil, err
	}

	return client, nil
}
