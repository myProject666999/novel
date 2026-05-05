package models

import (
	"time"

	"gorm.io/gorm"
)

type Novel struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title" gorm:"size:100;not null"`
	AuthorID     uint           `json:"author_id" gorm:"not null"`
	Author       User           `json:"author" gorm:"foreignKey:AuthorID"`
	CategoryID   uint           `json:"category_id" gorm:"not null"`
	Category     Category       `json:"category" gorm:"foreignKey:CategoryID"`
	Cover        string         `json:"cover" gorm:"size:255"`
	Description  string         `json:"description" gorm:"type:text"`
	Status       int            `json:"status" gorm:"default:1"` // 1:连载中, 2:已完结, 3:下架
	WordCount    int            `json:"word_count" gorm:"default:0"`
	ClickCount   int            `json:"click_count" gorm:"default:0"`
	CollectCount int            `json:"collect_count" gorm:"default:0"`
	CommentCount int            `json:"comment_count" gorm:"default:0"`
	Recommend    int            `json:"recommend" gorm:"default:0"` // 是否推荐: 0否, 1是
	VIP          int            `json:"vip" gorm:"default:0"` // 是否VIP: 0否, 1是
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Chapter struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	NovelID     uint           `json:"novel_id" gorm:"not null;index"`
	Novel       Novel          `json:"novel" gorm:"foreignKey:NovelID"`
	Title       string         `json:"title" gorm:"size:100;not null"`
	Content     string         `json:"content" gorm:"type:longtext"`
	WordCount   int            `json:"word_count" gorm:"default:0"`
	ChapterNum  int            `json:"chapter_num" gorm:"not null"`
	VIP         int            `json:"vip" gorm:"default:0"` // 是否VIP章节: 0否, 1是
	Price       float64        `json:"price" gorm:"default:0"`
	Status      int            `json:"status" gorm:"default:1"` // 1:正常, 2:下架
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Description string         `json:"description" gorm:"size:255"`
	Sort        int            `json:"sort" gorm:"default:0"`
	Status      int            `json:"status" gorm:"default:1"` // 1:启用, 2:禁用
	Novels      []Novel        `json:"novels,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	NovelID   uint           `json:"novel_id" gorm:"not null;index"`
	Novel     Novel          `json:"novel" gorm:"foreignKey:NovelID"`
	ParentID  *uint          `json:"parent_id"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Status    int            `json:"status" gorm:"default:1"` // 1:正常, 2:隐藏
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Bookshelf struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;uniqueIndex:idx_user_novel"`
	NovelID   uint           `json:"novel_id" gorm:"not null;uniqueIndex:idx_user_novel"`
	Novel     Novel          `json:"novel" gorm:"foreignKey:NovelID"`
	LastRead  *uint          `json:"last_read"` // 最后阅读的章节ID
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type NovelRecommend struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	NovelID   uint           `json:"novel_id" gorm:"not null"`
	Novel     Novel          `json:"novel" gorm:"foreignKey:NovelID"`
	Position  int            `json:"position" gorm:"default:0"` // 推荐位置
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"` // 1:启用, 2:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
