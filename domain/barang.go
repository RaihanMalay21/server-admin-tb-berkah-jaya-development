package domain

import (
	"gorm.io/gorm"
	"time"
)

type Barang struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
	Nama_Barang  string  `gorm:"Varchar(300); not null; unique"`
	Harga_Barang  float64  `gorm:"type:DECIMAL(10,0);not null"`
	Harga_Beli float64 `gorm:"type:DECIMAL(10,0);not null"`
	Image string `gorm:"varchar(200)"`
	Kode string  `gorm:"type:varchar(200); unique"`
}