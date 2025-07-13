package init

import (
	"log"

	"blog/internal/global"
)

// Initialize 初始化所有组件
func Initialize() error {
	// 1. 初始化配置
	if err := InitConfig(); err != nil {
		log.Fatal("Failed to initialize config:", err)
		return err
	}

	// 2. 初始化缓存
	global.InitCache()
	log.Println("Cache initialized successfully")

	// 3. 初始化数据库
	if err := InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
		return err
	}

	// 4. 初始化Redis
	if err := InitRedis(); err != nil {
		log.Fatal("Failed to initialize Redis:", err)
		return err
	}

	log.Println("All components initialized successfully")
	return nil
}
