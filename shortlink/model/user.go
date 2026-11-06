package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement",json:"id"`
	Name      string    `gorm:"uniqueIndex;size:64;not null",json:"username"`
	Password  string    `gorm:"size:255;not null",json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
