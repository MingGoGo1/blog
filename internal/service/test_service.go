package service

import (
	"blog/internal/global"
	"blog/model/entity"
	"blog/model/req"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TestService struct {
}

func NewTestService() *TestService {
	return &TestService{}
}

// getDB 获取数据库连接，支持事务
func (s *TestService) getDB(ctx context.Context) *gorm.DB {
	return global.GetDB(ctx)
}

// 创建test记录
func (s *TestService) CreateTest(ctx context.Context, req *req.TestCreateRequest) (*entity.Test, error) {

	// 获取db
	db := s.getDB(ctx)

	// 创建test记录
	test := &entity.Test{
		ID:       req.Id,
		Test:     req.Test,
		AuthorID: 1,
	}

	// 保存test记录
	if err := db.Create(test).Error; err != nil {
		return nil, fmt.Errorf("添加test记录失败: %v", err)
	}

	// 预加载关联数据

	return test, nil
}

func (s *TestService) Delete(ctx context.Context, id uint) error {
	db := s.getDB(ctx)

	var test entity.Test
	if err := db.First(&test, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("test不存在")
		}
		return err
	}

	return db.Delete(&test).Error
}

func (s *TestService) Update(ctx context.Context, id uint, req *req.TestUpdateRequest) (*entity.Test, error) {
	db := s.getDB(ctx)

	var test entity.Test
	if err := db.First(&test, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("test不存在")
		}
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Id != 0 {
		updates["id"] = req.Id
	}
	if req.Test != "" {
		updates["test"] = req.Test
	}

	if err := db.Model(&test).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新查询更新后的记录
	if err := db.Preload("Author").First(&test, id).Error; err != nil {
		return nil, err
	}

	return &test, nil
}

func (s *TestService) Get(ctx context.Context, id uint) (*entity.Test, error) {
	db := s.getDB(ctx)

	var test entity.Test
	if err := db.Preload("Author").First(&test, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("test不存在")
		}
		return nil, err
	}
	return &test, nil
}

func (s *TestService) List(ctx context.Context, req *req.TestListRequest) ([]*entity.Test, error) {
	db := s.getDB(ctx)

	// 查询条件
	query := db
	if req.Test != "" {
		query = query.Where("test LIKE ?", "%"+req.Test+"%")
	}

	// 分页
	query = query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize)

	// 查询
	var tests []*entity.Test
	if err := query.Find(&tests).Error; err != nil {
		return nil, err
	}

	return tests, nil

}
