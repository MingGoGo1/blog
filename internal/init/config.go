package init

import (
	"log"

	"blog/internal/config"
	"blog/internal/global"
)

// InitConfig 初始化配置
func InitConfig() error {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return err
	}

	global.Config = cfg
	log.Println("Config initialized successfully")
	return nil
}
