package main

import (
	"log"
	"os"
)

func main() {
	// 简单的测试程序
	log.Println("Laima 项目测试")
	log.Println("环境变量检查:")
	log.Printf("LAIMA_HTTP_PORT: %s", os.Getenv("LAIMA_HTTP_PORT"))
	log.Printf("LAIMA_DB_HOST: %s", os.Getenv("LAIMA_DB_HOST"))
	log.Printf("LAIMA_REDIS_URL: %s", os.Getenv("LAIMA_REDIS_URL"))
	log.Printf("LAIMA_MINIO_ENDPOINT: %s", os.Getenv("LAIMA_MINIO_ENDPOINT"))
	log.Printf("LAIMA_MEILISEARCH_URL: %s", os.Getenv("LAIMA_MEILISEARCH_URL"))
	log.Println("项目构建成功！")
}
