package domain

import (
	"gorm.io/gorm"
	"time"
)

type Hadiah struct {
	gorm.Model
	ID  uint `gorm:"id"` 
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
	Nama_Barang string `gorm:"type:varchar(200);not null"`
	Harga_Hadiah float64 `gorm:"type:DECIMAL(10,0);not null"`
	Poin float64 `gorm:"type:DECIMAL(10,0);not null"`
	Image string `gorm:"type:varchar(200); not null"`
	Deskripsi string `gorm:"type:varchar(500)"`
	User []User `gorm:"many2many:hadiah_users"`
}

type GetHadiah struct{
	Nama_Barang string `gorm:type:varchar(100);not null" json:"nama_barang`
	Poin int64
	Image string
	Desc string
}