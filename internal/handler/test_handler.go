package handler

import (
	"context"
	"strconv"

	"blog/internal/service"
	"blog/internal/utils"
	"blog/model/entity"
	"blog/model/req"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	testService *service.TestService
}

func NewTestHandler(testService *service.TestService) *TestHandler {
	return &TestHandler{
		testService: testService,
	}
}

// CreateTest 创建测试
// @Summary 创建测试
// @Description 创建新测试
// @Tags 测试管理
// @Accept json
// @Produce json
// @Param request body req.TestCreateRequest true "测试信息"
// @Success 200 {object} entity.Test "创建成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Router /api/v1/test [post]
func (h *TestHandler) CreateTest(c *gin.Context) {
	// 创建一个TestCreateRequest结构体实例
	var request req.TestCreateRequest

	// 绑定JSON请求数据到request结构体
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 使用统一事务处理
	test, err := utils.WithTransactionResult(c, func(ctx context.Context) (*entity.Test, error) {
		return h.testService.CreateTest(ctx, &request)
	})

	if err != nil {
		utils.Error(c, 2001, err.Error())
		return
	}

	utils.Success(c, test)
}

// DeleteTest 删除测试
// @Summary 删除测试
// @Description 删除指定测试
// @Tags 测试管理
// @Accept json
// @Produce json
// @Param id path int true "测试ID"
// @Success 200 {object} utils.Response "删除成功"
// @Failure 404 {object} utils.Response "测试不存在"
// @Router /api/v1/test/{id} [delete]
func (h *TestHandler) DeleteTest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID参数")
		return
	}

	// 使用统一事务处理
	err = utils.WithTransaction(c, func(ctx context.Context) error {
		return h.testService.Delete(ctx, uint(id))
	})

	if err != nil {
		utils.Error(c, 2002, err.Error())
		return
	}

	utils.Success(c, "删除成功")
}

// UpdateTest 更新测试
// @Summary 更新测试
// @Description 更新测试信息
// @Tags 测试管理
// @Accept json
// @Produce json
// @Param id path int true "测试ID"
// @Param request body req.TestUpdateRequest true "测试更新信息"
// @Success 200 {object} entity.Test "更新成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 404 {object} utils.Response "测试不存在"
// @Router /api/v1/test/{id} [put]
func (h *TestHandler) UpdateTest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID参数")
		return
	}

	var request req.TestUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 使用统一事务处理
	test, err := utils.WithTransactionResult(c, func(ctx context.Context) (*entity.Test, error) {
		return h.testService.Update(ctx, uint(id), &request)
	})

	if err != nil {
		utils.Error(c, 2003, err.Error())
		return
	}

	utils.Success(c, test)
}

// GetTest 获取测试详情
// @Summary 获取测试详情
// @Description 根据ID获取测试详细信息
// @Tags 测试管理
// @Accept json
// @Produce json
// @Param id path int true "测试ID"
// @Success 200 {object} entity.Test "获取成功"
// @Failure 404 {object} utils.Response "测试不存在"
// @Router /api/v1/test/{id} [get]
func (h *TestHandler) GetTest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID参数")
		return
	}

	// 调用服务方法
	test, err := h.testService.Get(c.Request.Context(), uint(id))
	if err != nil {
		utils.Error(c, 2004, err.Error())
		return
	}

	utils.Success(c, test)
}

// GetTests 获取测试列表
// @Summary 获取测试列表
// @Description 分页获取测试列表
// @Tags 测试管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} resp.TestListResponse "获取成功"
// @Router /api/v1/tests [get]
func (h *TestHandler) GetTests(c *gin.Context) {
	var request req.TestListRequest

	// 设置默认值
	request.Page = 1
	request.PageSize = 10

	// 绑定查询参数到request结构体
	if err := c.ShouldBindQuery(&request); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 调用服务方法
	tests, err := h.testService.List(c.Request.Context(), &request)
	if err != nil {
		utils.Error(c, 2005, err.Error())
		return
	}

	utils.Success(c, tests)
}
