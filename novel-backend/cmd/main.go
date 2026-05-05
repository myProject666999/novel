package main

import (
	"log"
	"novel-backend/config"
	"novel-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 创建Gin引擎
	r := gin.Default()

	// 配置路由
	routes.SetupRoutes(r)

	// 启动服务器
	addr := ":" + string(rune(config.AppConfig.Server.Port))
	log.Printf("Server starting on port %d", config.AppConfig.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
