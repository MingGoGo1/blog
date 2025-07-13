package entity

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID           uint           `json:"id" gorm:"primarykey;comment:文件ID"`
	FileName     string         `json:"file_name" gorm:"not null;size:255;comment:存储文件名"`
	OriginalName string         `json:"original_name" gorm:"not null;size:255;comment:原始文件名"`
	FileSize     int64          `json:"file_size" gorm:"not null;comment:文件大小(字节)"`
	FileType     string         `json:"file_type" gorm:"size:50;comment:文件类型/MIME类型"`
	FilePath     string         `json:"file_path" gorm:"not null;size:500;comment:文件存储路径"`
	FileURL      string         `json:"file_url" gorm:"not null;size:500;comment:文件访问URL"`
	UploaderID   uint           `json:"uploader_id" gorm:"not null;index;comment:上传者ID"`
	Status       int            `json:"status" gorm:"default:1;comment:状态:1-正常,0-已删除"`
	CreatedAt    time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联
	Uploader *User `json:"uploader,omitempty" gorm:"foreignKey:UploaderID"`
}
