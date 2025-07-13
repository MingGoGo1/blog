package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey;comment:用户ID"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:50;comment:用户名"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:100;comment:邮箱地址"`
	Password  string         `json:"-" gorm:"not null;size:255;comment:密码(加密)"`
	Nickname  string         `json:"nickname" gorm:"size:50;comment:昵称"`
	Avatar    string         `json:"avatar" gorm:"size:255;comment:头像URL"`
	Bio       string         `json:"bio" gorm:"size:500;comment:个人简介"`
	Status    int            `json:"status" gorm:"default:1;comment:状态:1-激活,0-禁用"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联
	Articles []Article `json:"articles,omitempty" gorm:"foreignKey:AuthorID"`
}
