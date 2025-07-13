package handler

import (
	"context"

	"blog/internal/middleware"
	"blog/internal/service"
	"blog/internal/utils"
	"blog/model/entity"
	"blog/model/req"
	"blog/model/resp"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body req.UserRegisterRequest true "注册信息"
// @Success 200 {object} entity.User "注册成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Router /api/v1/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var request req.UserRegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 使用统一事务处理
	user, err := utils.WithTransactionResult(c, func(ctx context.Context) (*entity.User, error) {
		return h.userService.Register(ctx, &request)
	})
	if err != nil {
		utils.Error(c, 1001, err.Error())
		return
	}

	// 隐藏密码
	user.Password = ""
	utils.Success(c, user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取JWT token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body req.UserLoginRequest true "登录信息"
// @Success 200 {object} resp.UserLoginResponse "登录成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "用户名或密码错误"
// @Router /api/v1/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var request req.UserLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 使用统一事务处理
	resp, err := utils.WithTransactionResult(c, func(ctx context.Context) (*resp.UserLoginResponse, error) {
		return h.userService.Login(ctx, &request)
	})
	if err != nil {
		utils.Error(c, 1002, err.Error())
		return
	}

	// 隐藏密码
	resp.User.Password = ""
	utils.Success(c, resp)
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户的详细资料
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.UserProfileResponse "获取成功"
// @Failure 401 {object} utils.Response "未登录"
// @Router /api/v1/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		utils.Error(c, 1003, "获取用户信息失败")
		return
	}

	utils.Success(c, resp.ToUserProfileResponse(user))
}

// Logout 用户注销
// @Summary 用户注销
// @Description 注销当前用户，清除token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "注销成功"
// @Failure 401 {object} utils.Response "未登录"
// @Router /api/v1/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	// 获取token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.Unauthorized(c, "缺少认证token")
		return
	}

	// 提取token
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		utils.Unauthorized(c, "认证token格式错误")
		return
	}

	// 注销token
	err := h.userService.Logout(c.Request.Context(), token)
	if err != nil {
		utils.Error(c, 1004, "注销失败: "+err.Error())
		return
	}

	utils.Success(c, "注销成功")
}
