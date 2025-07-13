package global

import (
	"context"
	"gorm.io/gorm"
)

// DB 全局数据库变量
var DB *gorm.DB

// GetDB 获取数据库连接，支持事务
func GetDB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return DB
	}

	// 从上下文中获取事务
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
		return tx
	}

	return DB
}
