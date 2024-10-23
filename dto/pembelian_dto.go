package dto

import (
	"time"
)

type Pembelian struct {
	ID uint `json:"ID" validate:"required"`
	UserID uint `json:"userid" validate:"required"`
	User User `validate:"-"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Tanggal_Pembelian  string `json:"tanggal_pembelian" validate:"required"`
	Total_Harga float64 `json:"total_harga" validate:"required"`
	Total_Keuntungan float64 `json:"total_keuntungan" validate:"required"`
	Image string `json:"image" validate:"required"`
	KeteranganNotaCancelID uint `json:"kerangan_nota_cancel_id"`
	Status string 
	KeteranganNotaCancel KeteranganNotaCancel 
}

type KeteranganNotaCancel struct {
	ID uint `json:"ID"`
	Desc string `json:"desc"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type PembelianResponse struct {
	ID               uint      `json:"ID"`
	CreatedAt        time.Time `json:"CreatedAt"`
	UpdatedAt        time.Time `json:"UpdatedAt"`
	UserID           uint      `json:"userid"`
	UserName		 string    `json:"username"`
	Email			 string    `json:"email"`
	Tanggal_Pembelian  string  `json:"tanggal_pembelian"`
	Total_Harga      float64   `json:"total_harga"`
	Total_Keuntungan float64   `json:"total_keuntungan"`
	Image            string    `json:"image"`
	Keterangan		string	   `json:"keterangan"`
}
