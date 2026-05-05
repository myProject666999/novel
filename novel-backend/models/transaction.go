package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	OrderNo        string         `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	UserID         uint           `json:"user_id" gorm:"not null"`
	User           User           `json:"user" gorm:"foreignKey:UserID"`
	NovelID        *uint          `json:"novel_id"`
	ChapterID      *uint          `json:"chapter_id"`
	OrderType      int            `json:"order_type" gorm:"not null"` // 1:充值, 2:订阅VIP, 3:购买章节
	Amount         float64        `json:"amount" gorm:"not null"`
	Status         int            `json:"status" gorm:"default:0"` // 0:待支付, 1:已支付, 2:已取消, 3:已退款
	PayTime        *time.Time     `json:"pay_time"`
	PayMethod      string         `json:"pay_method" gorm:"size:20"`
	TransactionID  string         `json:"transaction_id" gorm:"size:100"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type RechargeRecord struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	OrderID       uint           `json:"order_id" gorm:"not null"`
	Amount        float64        `json:"amount" gorm:"not null"`
	BalanceBefore float64        `json:"balance_before" gorm:"not null"`
	BalanceAfter  float64        `json:"balance_after" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type Subscription struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;uniqueIndex:idx_user_novel"`
	NovelID   uint           `json:"novel_id" gorm:"not null;uniqueIndex:idx_user_novel"`
	Novel     Novel          `json:"novel" gorm:"foreignKey:NovelID"`
	Status    int            `json:"status" gorm:"default:1"` // 1:订阅中, 2:已取消
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Member struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;uniqueIndex"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Level       int            `json:"level" gorm:"default:1"` // 会员等级
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	Status      int            `json:"status" gorm:"default:1"` // 1:有效, 2:过期
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Feedback struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Title     string         `json:"title" gorm:"size:100;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Status    int            `json:"status" gorm:"default:0"` // 0:待处理, 1:已处理
	Reply     string         `json:"reply" gorm:"type:text"`
	ReplyTime *time.Time     `json:"reply_time"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type SiteInfo struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	SiteName    string         `json:"site_name" gorm:"size:100"`
	SiteLogo    string         `json:"site_logo" gorm:"size:255"`
	SiteKeywords string        `json:"site_keywords" gorm:"size:255"`
	SiteDescription string     `json:"site_description" gorm:"size:500"`
	Copyright   string         `json:"copyright" gorm:"size:255"`
	Icp         string         `json:"icp" gorm:"size:50"`
	ContactEmail string        `json:"contact_email" gorm:"size:100"`
	ContactPhone string        `json:"contact_phone" gorm:"size:20"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type FriendLink struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:50;not null"`
	Url       string         `json:"url" gorm:"size:255;not null"`
	Logo      string         `json:"logo" gorm:"size:255"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"` // 1:启用, 2:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type News struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:100;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Cover     string         `json:"cover" gorm:"size:255"`
	Author    string         `json:"author" gorm:"size:50"`
	Views     int            `json:"views" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"` // 1:发布, 2:草稿
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type SystemLog struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id"`
	Username  string         `json:"username" gorm:"size:50"`
	Module    string         `json:"module" gorm:"size:50"`
	Action    string         `json:"action" gorm:"size:100"`
	Method    string         `json:"method" gorm:"size:10"`
	Path      string         `json:"path" gorm:"size:255"`
	IP        string         `json:"ip" gorm:"size:50"`
	Params    string         `json:"params" gorm:"type:text"`
	Result    string         `json:"result" gorm:"type:text"`
	Status    int            `json:"status"` // 1:成功, 2:失败
	CreatedAt time.Time      `json:"created_at"`
}
