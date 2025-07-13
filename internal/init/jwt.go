package init

import (
	"fmt"
	"log"

	"blog/internal/global"
	"blog/internal/utils"
)

// InitJWT 初始化JWT配置
func InitJWT() error {
	cfg := global.Config
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	// 设置JWT密钥
	utils.SetJWTSecret(cfg.JWT.Secret)
	
	// 设置JWT过期时间
	utils.SetJWTExpireHour(cfg.JWT.ExpireHour)

	log.Println("JWT initialized successfully")
	return nil
}
