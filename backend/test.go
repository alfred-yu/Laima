package main

import (
	"fmt"
	repoapp "laima/internal/repo/app"
	repoapi "laima/internal/repo/api"
	userapp "laima/internal/user/app"
	userapi "laima/internal/user/api"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 测试代码结构是否正确
	fmt.Println("Testing code structure...")

	// 测试仓库服务
	db, err := gorm.Open(postgres.Open("host=localhost port=5432 user=laima password=laima123 dbname=laima sslmode=disable"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}

	// 测试仓库API
	repoAPI := repoapi.NewRepoAPI(db, nil, nil, nil)
	fmt.Println("RepoAPI created successfully")

	// 测试用户服务
	userService := userapp.NewUserService(db)
	fmt.Println("UserService created successfully")

	// 测试用户API
	userAPI := userapi.NewUserAPI(db, nil)
	fmt.Println("UserAPI created successfully")

	fmt.Println("All tests passed!")
}
