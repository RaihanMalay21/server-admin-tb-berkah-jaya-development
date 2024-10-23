package domain

import (
	"gorm.io/gorm"
	"time"
)

type Pembelian struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`
	UserID uint 
	User User `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
	Tanggal_Pembelian  string `gorm:"type:varchar(200)"`
	Total_Harga float64 `gorm:"type:DECIMAL(10, 0)"`
	Total_Keuntungan float64 `gorm:"type:DECIMAL(10, 0)"`
	Image string `gorm:"type:varchar(300); not null"`
	KeteranganNotaCancelID uint 
	Status string `gorm:"type:enum('cancel', 'success');default:'cancel'"`
	KeteranganNotaCancel KeteranganNotaCancel `gorm:"foreignKey:KeteranganNotaCancelID;references:ID"`
}

type KeteranganNotaCancel struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`
	Desc string `gorm:"type:varchar(200)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
}


