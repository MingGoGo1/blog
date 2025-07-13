package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"blog/internal/global"
	"blog/model/entity"
	"blog/model/resp"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

type FileService struct {
	minioClient *minio.Client
}

func NewFileService() *FileService {
	// 初始化MinIO客户端
	minioClient, err := InitMinio()
	if err != nil {
		// 如果MinIO初始化失败，记录错误但不终止程序
		fmt.Printf("Warning: Failed to initialize MinIO: %v\n", err)
		minioClient = nil
	}

	return &FileService{
		minioClient: minioClient,
	}
}

// getDB 获取数据库连接，支持事务
func (s *FileService) getDB(ctx context.Context) *gorm.DB {
	return global.GetDB(ctx)
}

// InitMinio 初始化MinIO客户端
func InitMinio() (*minio.Client, error) {
	cfg := global.Config
	if cfg == nil {
		return nil, fmt.Errorf("config not initialized")
	}

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// 检查并创建bucket
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, cfg.Minio.BucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.Minio.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}

		// 设置bucket为公共读取权限
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + cfg.Minio.BucketName + `/*"]
				}
			]
		}`

		err = minioClient.SetBucketPolicy(ctx, cfg.Minio.BucketName, policy)
		if err != nil {
			fmt.Printf("Warning: Failed to set bucket policy: %v\n", err)
			// 不返回错误，因为bucket已经创建成功
		}
	} else {
		// 如果bucket已存在，检查并设置策略
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + cfg.Minio.BucketName + `/*"]
				}
			]
		}`

		err = minioClient.SetBucketPolicy(ctx, cfg.Minio.BucketName, policy)
		if err != nil {
			fmt.Printf("Warning: Failed to set bucket policy for existing bucket: %v\n", err)
		}
	}

	return minioClient, nil
}

// SetBucketPublicPolicy 设置bucket为公共读取权限
// 可以用于修复已存在的bucket权限问题
func SetBucketPublicPolicy() error {
	cfg := global.Config
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
	if err != nil {
		return err
	}

	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::` + cfg.Minio.BucketName + `/*"]
			}
		]
	}`

	ctx := context.Background()
	err = minioClient.SetBucketPolicy(ctx, cfg.Minio.BucketName, policy)
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %v", err)
	}

	fmt.Printf("Successfully set public read policy for bucket: %s\n", cfg.Minio.BucketName)
	return nil
}

// UploadFile 上传文件
func (s *FileService) UploadFile(ctx context.Context, userID uint, fileHeader *multipart.FileHeader) (*resp.FileUploadResponse, error) {
	// 检查MinIO客户端是否可用
	if s.minioClient == nil {
		return nil, fmt.Errorf("文件上传服务不可用，MinIO未正确配置")
	}

	db := s.getDB(ctx)

	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 生成文件名
	ext := filepath.Ext(fileHeader.Filename)
	dateStr := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
	dirPath := fmt.Sprintf("uploads/%s", dateStr)
	filePath := fmt.Sprintf("%s/%s", dirPath, fileName)

	// 确保目录存在
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return nil, err
	}

	// 获取文件类型
	fileType := getFileType(ext)

	// 上传到MinIO
	cfg := global.Config
	_, err = s.minioClient.PutObject(ctx, cfg.Minio.BucketName, filePath, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: fileType,
	})
	if err != nil {
		return nil, err
	}

	// 生成文件URL
	fileURL := s.generateFileURL(filePath)

	// 保存文件记录到数据库
	fileRecord := &entity.File{
		FileName:     fileName,
		OriginalName: fileHeader.Filename,
		FileSize:     fileHeader.Size,
		FileType:     fileType,
		FilePath:     filePath,
		FileURL:      fileURL,
		UploaderID:   userID,
		Status:       1,
	}

	if err := db.Create(fileRecord).Error; err != nil {
		return nil, err
	}

	return &resp.FileUploadResponse{
		ID:       fileRecord.ID,
		FileName: fileRecord.FileName,
		FileURL:  fileRecord.FileURL,
		FileSize: fileRecord.FileSize,
		FileType: fileRecord.FileType,
	}, nil
}

// GetFileByID 根据ID获取文件信息
func (s *FileService) GetFileByID(ctx context.Context, id uint) (*entity.File, error) {
	db := s.getDB(ctx)

	var file entity.File
	if err := db.First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(ctx context.Context, id, userID uint) error {
	db := s.getDB(ctx)

	var file entity.File
	if err := db.First(&file, id).Error; err != nil {
		return err
	}

	// 检查权限
	if file.UploaderID != userID {
		return fmt.Errorf("无权限删除此文件")
	}

	// 从MinIO删除文件
	cfg := global.Config
	err := s.minioClient.RemoveObject(ctx, cfg.Minio.BucketName, file.FilePath, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	// 从数据库删除记录
	return db.Delete(&file).Error
}

// generateFileURL 生成文件访问URL
func (s *FileService) generateFileURL(filePath string) string {
	cfg := global.Config
	protocol := "http"
	if cfg.Minio.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, cfg.Minio.Endpoint, cfg.Minio.BucketName, filePath)
}

// getFileType 根据文件扩展名获取文件类型
func getFileType(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}
