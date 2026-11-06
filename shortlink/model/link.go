package model

import "time"

type Link struct {
	ID        uint      `gorm:"primaryKey;autoIncrement",json:"id"`
	UserID    uint      `gorm:"index;not null"json:"user_id"`
	Code      string    `gorm:"size:32;uniqueIndex;not null"json:"code"`
	LongUrl   string    `gorm:"size:2048;not null"json:"long_url"`
	Clicks    int64     `gorm:"default:0"json:"clicks"`
	CreatedAt time.Time `gorm:"autoCreateTime"json:"created_at"`
}
