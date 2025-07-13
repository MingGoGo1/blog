package service

import (
	"context"
	"errors"
	"time"

	"blog/internal/global"
	"blog/model/entity"
	"blog/model/req"
	"blog/model/resp"

	"gorm.io/gorm"
)

type ArticleService struct{}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

// getDB 获取数据库连接，支持事务
func (s *ArticleService) getDB(ctx context.Context) *gorm.DB {
	return global.GetDB(ctx)
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(ctx context.Context, userID uint, req *req.ArticleCreateRequest) (*entity.Article, error) {
	db := s.getDB(ctx)

	article := &entity.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		CoverImage: req.CoverImage,
		AuthorID:   userID,
		Status:     req.Status,
		Tags:       req.Tags,
	}

	// 如果是发布状态，设置发布时间
	if req.Status == 1 {
		now := time.Now()
		article.PublishedAt = &now
	}

	if err := db.Create(article).Error; err != nil {
		return nil, err
	}

	// 预加载作者信息
	if err := db.Preload("Author").First(article, article.ID).Error; err != nil {
		return nil, err
	}

	return article, nil
}

// GetArticles 获取文章列表
func (s *ArticleService) GetArticles(ctx context.Context, req *req.ArticleListRequest) (*resp.ArticleListResponse, error) {
	db := s.getDB(ctx)

	var articles []entity.Article
	var total int64

	query := db.Model(&entity.Article{})

	// 添加查询条件
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}
	if req.AuthorID != 0 {
		query = query.Where("author_id = ?", req.AuthorID)
	}
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Author").
		Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&articles).Error; err != nil {
		return nil, err
	}

	// 处理响应数据
	for i := range articles {
		resp.ToArticleResponse(&articles[i])
	}

	return &resp.ArticleListResponse{
		Articles: articles,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetArticleByID 根据ID获取文章
func (s *ArticleService) GetArticleByID(ctx context.Context, id uint) (*entity.Article, error) {
	db := s.getDB(ctx)

	var article entity.Article
	if err := db.Preload("Author").First(&article, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	// 增加浏览量
	db.Model(&article).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	return resp.ToArticleResponse(&article), nil
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(ctx context.Context, id, userID uint, req *req.ArticleUpdateRequest) (*entity.Article, error) {
	db := s.getDB(ctx)

	var article entity.Article
	if err := db.First(&article, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	// 检查权限
	if article.AuthorID != userID {
		return nil, errors.New("无权限修改此文章")
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.CoverImage != "" {
		updates["cover_image"] = req.CoverImage
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}

	// 处理状态变更
	if req.Status != article.Status {
		updates["status"] = req.Status
		if req.Status == 1 && article.PublishedAt == nil {
			// 首次发布
			now := time.Now()
			updates["published_at"] = &now
		}
	}

	if err := db.Model(&article).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新查询更新后的文章
	if err := db.Preload("Author").First(&article, id).Error; err != nil {
		return nil, err
	}

	return resp.ToArticleResponse(&article), nil
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(ctx context.Context, id, userID uint) error {
	db := s.getDB(ctx)

	var article entity.Article
	if err := db.First(&article, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}

	// 检查权限
	if article.AuthorID != userID {
		return errors.New("无权限删除此文章")
	}

	return db.Delete(&article).Error
}
