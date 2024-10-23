package dto

import (
	"time"
)

type Pembelian_Per_Item struct {
	ID uint `json:"ID"`
	PembelianID uint `json:"pembelianID" validate:"required"`
	BarangID uint `json:"barangID" validate:"required"`
	Barang Barang `validate:"-"`
	// Pembelian Pembelian `gorm:"foreignKey:PembelianID"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Jumlah_Barang float64 `json:"jumlah_barang" validate:"required"`
	Total_Harga float64 `json:"total_harga" validate:"required"`
	Total_Keuntungan float64 `json:"total_keuntungan" validate:"required"`
}