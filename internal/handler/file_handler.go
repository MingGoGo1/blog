package handler

import (
	"context"

	"blog/internal/middleware"
	"blog/internal/service"
	"blog/internal/utils"
	"blog/model/resp"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// Upload 上传文件
// @Summary 上传文件
// @Description 上传文件到服务器
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "要上传的文件"
// @Success 200 {object} resp.FileUploadResponse "上传成功"
// @Failure 400 {object} utils.Response "文件格式错误或文件过大"
// @Failure 401 {object} utils.Response "未登录"
// @Router /api/v1/upload [post]
func (h *FileHandler) Upload(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.Unauthorized(c, "未登录")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 检查文件大小（限制为10MB）
	if file.Size > 10*1024*1024 {
		utils.BadRequest(c, "文件大小不能超过10MB")
		return
	}

	// 使用统一事务处理上传文件
	resp, err := utils.WithTransactionResult(c, func(ctx context.Context) (*resp.FileUploadResponse, error) {
		return h.fileService.UploadFile(ctx, userID, file)
	})
	if err != nil {
		utils.Error(c, 3001, "文件上传失败: "+err.Error())
		return
	}

	utils.Success(c, resp)
}
