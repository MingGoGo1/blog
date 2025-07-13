package init

import (
	"fmt"
	"log"
	"time"

	"blog/internal/global"
	"blog/model/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB() error {
	cfg := global.Config
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	// 首先连接MySQL服务器（不指定数据库）
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	// 连接MySQL服务器
	tempDB, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 临时连接不显示日志
	})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %v", err)
	}

	// 创建数据库（如果不存在）
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Database.DBName)
	if err := tempDB.Exec(createDBSQL).Error; err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}

	// 关闭临时连接
	tempSqlDB, _ := tempDB.DB()
	tempSqlDB.Close()

	// 现在连接到指定的数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移
	if err := autoMigrate(db); err != nil {
		return err
	}

	global.DB = db
	log.Println("Database initialized successfully")
	return nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Article{},
		&entity.File{},
		&entity.Test{},
	)
}
