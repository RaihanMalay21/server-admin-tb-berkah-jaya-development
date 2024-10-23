package dto

import (
	"time"
)

type HadiahUser struct {
	UserID  uint 
	HadiahID uint 
	Hadiah Hadiah 
	User User 
	GiftsArrive string 
	Status string 
	CreatedAt time.Time
	UpdatedAt time.Time
}