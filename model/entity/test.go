package entity

import (
	"gorm.io/gorm"
	"time"
)

type Test struct {
	ID          uint           `json:"id" gorm:"primarykey;comment:主键"`
	Test        string         `json:"test" gorm:"not null;size:200;comment:测试内容"`
	AuthorID    uint           `json:"author_id" gorm:"not null;index;comment:测试关联id"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	PublishedAt *time.Time     `json:"published_at" gorm:"comment:发布时间"`

	// 关联
	Author *User `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
}
