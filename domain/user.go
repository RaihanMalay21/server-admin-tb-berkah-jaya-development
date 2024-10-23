package domain

import( 
"gorm.io/gorm"
"time"
)

type User struct {
	gorm.Model
	ID uint `gorm:"primarykey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"` 
	UserName string `gorm:"type:varchar(100);unique;not null"`
	Email string	`gorm:"type:varchar(200);unique;not null"`
	NoWhatshapp string `gorm:"type:varchar(20);unique;not null"`
	Password string `gorm:"type:varchar(200);not null"` 
	Poin float64   `gorm:"type:DECIMAL(10, 0)"`
	Pembelian []Pembelian 
	Hadiah  []Hadiah `gorm:"many2many:hadiah_users"`
}