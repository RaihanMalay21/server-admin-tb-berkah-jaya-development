package domain

import (
	"gorm.io/gorm"
	"time"
)

type HadiahUser struct {
	gorm.Model
	UserID  uint `gorm:"primaryKey;not null"`
	HadiahID uint `gorm:"primaryKey; not null"`
	Hadiah Hadiah `gorm:"foreignKey:HadiahID;references:ID"` 
	User User `gorm:"foreignKey:UserID;references:ID"`
	GiftsArrive string `gorm:"type:enum('NO','YES');default:'NO'"`
	Status string `gorm:"type:enum('unfinished', 'finished');default:'unfinished'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}