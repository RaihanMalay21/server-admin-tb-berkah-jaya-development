package dto

import (
	"time"
)

type Hadiah struct {
	ID  uint `json:"ID"` 
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Nama_Barang string `json:"nama_barang" validate:"required"`
	Harga_Hadiah float64 `json:"harga_hadiah" validate:"required,numeric"`
	Poin float64 `json:"poin" validate:"required,numeric"`
	Image string `json:"image"`
	Deskripsi string `json:"desc" validate:"max=500"`
	User []User 
}

type GetHadiah struct{
	Nama_Barang string `json:"nama_barang`
	Poin int64
	Image string
	Desc string
}