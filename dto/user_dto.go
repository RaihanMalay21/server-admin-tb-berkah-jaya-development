package dto

import ( 
	"time"
	)
	
type User struct {
	ID uint `json:"id"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"` 
	UserName string `json:"username" validate:"required,alphanum,uniqueUsername"`
	Email string	`json:"email" validate:"required,email,uniqueEmail"`
	NoWhatshapp string `json:"whatshapp" validate:"required,uniquePhone"`
	Password string `json:"password" validate:"required,alphanum,min=6"` 
	Poin float64   `json:"poin"`
	Pembelian []Pembelian 
	Hadiah  []Hadiah 
}