package handler

import (
	"context"
	"strconv"

	"blog/internal/middleware"
	"blog/internal/service"
	"blog/internal/utils"
	"blog/model/entity"
	"blog/model/req"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService *service.ArticleService
}

func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

// CreateArticle 创建文章
// @Summary 创建文章
// @Description 创建新文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body req.ArticleCreateRequest true "文章信息"
// @Success 200 {object} entity.Article "创建成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未登录"
// @Router /api/v1/articles [post]
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	var request req.ArticleCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 使用统一事务处理
	article, err := utils.WithTransactionResult(c, func(ctx context.Context) (*entity.Article, error) {
		return h.articleService.CreateArticle(ctx, userID, &request)
	})
	if err != nil {
		utils.Error(c, 2001, err.Error())
		return
	}

	utils.Success(c, article)
}

// GetArticles 获取文章列表
// @Summary 获取文章列表
// @Description 分页获取文章列表，支持搜索和筛选
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param author_id query int false "作者ID筛选"
// @Param status query int false "文章状态" Enums(0, 1) default(1)
// @Success 200 {object} resp.ArticleListResponse "获取成功"
// @Router /api/v1/articles [get]
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	var request req.ArticleListRequest

	// 设置默认值
	request.Page = 1
	request.PageSize = 10
	request.Status = 1 // 默认只显示已发布的文章

	// 绑定查询参数
	if err := c.ShouldBindQuery(&request); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	resp, err := h.articleService.GetArticles(c.Request.Context(), &request)
	if err != nil {
		utils.Error(c, 2002, err.Error())
		return
	}

	utils.Success(c, resp)
}

// GetArticle 获取单篇文章
// @Summary 获取文章详情
// @Description 根据ID获取文章详细信息
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} entity.Article "获取成功"
// @Failure 404 {object} utils.Response "文章不存在"
// @Router /api/v1/articles/{id} [get]
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	article, err := h.articleService.GetArticleByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.Error(c, 2003, err.Error())
		return
	}

	utils.Success(c, article)
}

// UpdateArticle 更新文章
// @Summary 更新文章
// @Description 更新文章信息
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Param request body req.ArticleUpdateRequest true "文章更新信息"
// @Success 200 {object} entity.Article "更新成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未登录"
// @Failure 404 {object} utils.Response "文章不存在"
// @Router /api/v1/articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	var request req.ArticleUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 使用统一事务处理
	article, err := utils.WithTransactionResult(c, func(ctx context.Context) (*entity.Article, error) {
		return h.articleService.UpdateArticle(ctx, uint(id), userID, &request)
	})
	if err != nil {
		utils.Error(c, 2004, err.Error())
		return
	}

	utils.Success(c, article)
}

// DeleteArticle 删除文章
// @Summary 删除文章
// @Description 删除指定文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} utils.Response "删除成功"
// @Failure 401 {object} utils.Response "未登录"
// @Failure 404 {object} utils.Response "文章不存在"
// @Router /api/v1/articles/{id} [delete]
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	// 使用统一事务处理
	err = utils.WithTransaction(c, func(ctx context.Context) error {
		return h.articleService.DeleteArticle(ctx, uint(id), userID)
	})
	if err != nil {
		utils.Error(c, 2005, err.Error())
		return
	}

	utils.Success(c, nil)
}
