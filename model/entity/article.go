package entity

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID          uint           `json:"id" gorm:"primarykey;comment:文章ID"`
	Title       string         `json:"title" gorm:"not null;size:200;comment:文章标题"`
	Content     string         `json:"content" gorm:"type:longtext;comment:文章内容"`
	Summary     string         `json:"summary" gorm:"size:500;comment:文章摘要"`
	CoverImage  string         `json:"cover_image" gorm:"size:255;comment:封面图片URL"`
	AuthorID    uint           `json:"author_id" gorm:"not null;index;comment:作者ID"`
	Status      int            `json:"status" gorm:"default:1;comment:状态:1-已发布,0-草稿"`
	ViewCount   int            `json:"view_count" gorm:"default:0;comment:浏览次数"`
	LikeCount   int            `json:"like_count" gorm:"default:0;comment:点赞次数"`
	Tags        string         `json:"tags" gorm:"size:255;comment:标签,逗号分隔"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	PublishedAt *time.Time     `json:"published_at" gorm:"comment:发布时间"`

	// 关联
	Author *User `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
}
