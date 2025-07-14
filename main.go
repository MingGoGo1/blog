// @title Blog API
// @version 1.0
// @description 博客系统API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8868
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"blog/internal/global"
	initpkg "blog/internal/init"
	"log"
)

func main() {
	// 初始化所有组件
	if err := initpkg.Initialize(); err != nil {
		log.Fatal("Failed to initialize:", err)
	}

	// 初始化HTTP服务器
	r := initpkg.InitHTTPServer()

	// 启动服务器
	log.Printf("Server starting on port %s", global.Config.Server.Port)
	if err := r.Run(":" + global.Config.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
