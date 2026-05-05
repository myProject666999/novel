package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password     string         `json:"-" gorm:"size:255;not null"`
	Nickname     string         `json:"nickname" gorm:"size:50"`
	Email        string         `json:"email" gorm:"size:100"`
	Phone        string         `json:"phone" gorm:"size:20"`
	Avatar       string         `json:"avatar" gorm:"size:255"`
	Balance      float64        `json:"balance" gorm:"default:0"`
	VipLevel     int            `json:"vip_level" gorm:"default:0"`
	Status       int            `json:"status" gorm:"default:1"` // 1:正常, 2:禁用
	RoleID       uint           `json:"role_id" gorm:"not null"`
	Role         Role           `json:"role" gorm:"foreignKey:RoleID"`
	InviteCodeID *uint          `json:"invite_code_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:50;not null"`
	DisplayName string         `json:"display_name" gorm:"size:50;not null"`
	Description string         `json:"description" gorm:"size:255"`
	Status      int            `json:"status" gorm:"default:1"` // 1:启用, 2:禁用
	Users       []User         `json:"users,omitempty" gorm:"foreignKey:RoleID"`
	Menus       []Menu         `json:"menus,omitempty" gorm:"many2many:role_menus"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Menu struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ParentID    *uint          `json:"parent_id"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Path        string         `json:"path" gorm:"size:255"`
	Icon        string         `json:"icon" gorm:"size:100"`
	Sort        int            `json:"sort" gorm:"default:0"`
	Status      int            `json:"status" gorm:"default:1"` // 1:启用, 2:禁用
	MenuType    int            `json:"menu_type" gorm:"default:1"` // 1:目录, 2:菜单, 3:按钮
	Permission  string         `json:"permission" gorm:"size:100"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type InviteCode struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Code      string         `json:"code" gorm:"uniqueIndex;size:20;not null"`
	UsedCount int            `json:"used_count" gorm:"default:0"`
	MaxCount  int            `json:"max_count" gorm:"default:1"`
	Status    int            `json:"status" gorm:"default:1"` // 1:启用, 2:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
