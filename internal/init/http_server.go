package init

import (
	"log"

	"blog/internal/global"
	"blog/internal/router"

	"github.com/gin-gonic/gin"
)

// InitHTTPServer 初始化HTTP服务器
func InitHTTPServer() *gin.Engine {
	cfg := global.Config
	if cfg == nil {
		log.Fatal("Config not initialized")
	}

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := router.SetupRouter()

	log.Println("HTTP server initialized successfully")
	return r
}
