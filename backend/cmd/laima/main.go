package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	repoapi "laima/internal/repo/api"
	userapi "laima/internal/user/api"
	prapi "laima/internal/pr/api"
	aiapi "laima/internal/ai/api"
	cicdapi "laima/internal/cicd/api"
	issueapi "laima/internal/issue/api"
	"laima/internal/git"
	"laima/internal/user/domain"
	repodomain "laima/internal/repo/domain"
	prdomain "laima/internal/pr/domain"
	issuedomain "laima/internal/issue/domain"
	cicddomain "laima/internal/cicd/domain"
	aidomain "laima/internal/ai/domain"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/meilisearch/meilisearch-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库连接
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表
	err = autoMigrate(db)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化 Redis 连接
	redisClient, err := initRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 初始化 MinIO 客户端
	minioClient, err := initMinIO()
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	// 初始化 Meilisearch 客户端
	var meiliClient meilisearch.ServiceManager
	meiliClient, err = initMeilisearch()
	if err != nil {
		log.Fatalf("Failed to connect to Meilisearch: %v", err)
	}

	// 初始化 Git 服务
	gitSvc := initGitService()

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

	// 注册 API 路由
	repoAPI := repoapi.NewRepoAPI(db, redisClient, minioClient, meiliClient, gitSvc)
	repoAPI.RegisterRoutes(r)

	userAPI := userapi.NewUserAPI(db, redisClient)
	userAPI.RegisterRoutes(r)

	prAPI := prapi.NewPRAPI(db, gitSvc)
	prAPI.RegisterRoutes(r)

	aiAPI := aiapi.NewAIApi(db)
	aiAPI.RegisterRoutes(r)

	cicdAPI := cicdapi.NewCICDApi(db)
	cicdAPI.RegisterRoutes(r)

	issueAPI := issueapi.NewIssueApi(db)
	issueAPI.RegisterRoutes(r)

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