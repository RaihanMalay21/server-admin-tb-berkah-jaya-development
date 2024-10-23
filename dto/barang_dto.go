package dto 

import (
	"time"
)

type Barang struct {
	ID uint  `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at`
	Nama_Barang  string `json:"nama_barang" validate:"required"`
	Harga_Barang  float64  `json:"harga_barang" validate:"required,numeric"`
	Harga_Beli float64 `json:"harga_beli" validate:"required,numeric"`
	Image string `json:"image"`
	Kode string  `json:"kode" validate:"required"`
}