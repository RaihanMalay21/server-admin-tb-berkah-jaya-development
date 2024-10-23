package domain

import (
	"gorm.io/gorm"
	"time"
)

type Pembelian_Per_Item struct {
	gorm.Model
	ID uint `gorm:"PrimaryKey"`
	PembelianID uint 
	BarangID uint 
	Barang Barang `gorm:"foreignKey:BarangID; references:ID"`
	// Pembelian Pembelian `gorm:"foreignKey:PembelianID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
	Jumlah_Barang float64 `gorm:"type:DECIMAL(10, 0);not null"`
	Total_Harga float64 `gorm:"type:DECIMAL(10, 0);not null"`
	Total_Keuntungan float64 `gorm:"type:DECIMAL(10, 0);not null"`
}