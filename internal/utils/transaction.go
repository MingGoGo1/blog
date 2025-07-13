package utils

import (
	"context"

	"blog/internal/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TransactionFunc 事务处理函数类型
type TransactionFunc func(ctx context.Context) error

// WithTransaction 统一事务处理封装
// 在Handler层使用，自动处理事务的开始、提交和回滚
func WithTransaction(c *gin.Context, fn TransactionFunc) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 将事务放入上下文
		ctx := context.WithValue(c.Request.Context(), "tx", tx)

		// 执行业务逻辑
		return fn(ctx)
	})
}

// WithTransactionResult 带返回值的事务处理封装
// 用于需要返回结果的事务操作
func WithTransactionResult[T any](c *gin.Context, fn func(ctx context.Context) (T, error)) (T, error) {
	var result T
	var err error

	txErr := global.DB.Transaction(func(tx *gorm.DB) error {
		// 将事务放入上下文
		ctx := context.WithValue(c.Request.Context(), "tx", tx)

		// 执行业务逻辑
		result, err = fn(ctx)
		return err
	})

	if txErr != nil {
		var zero T
		return zero, txErr
	}

	return result, err
}
