package repository

import (
	"fmt"
	"errors"
	"gorm.io/gorm"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
)

type BarangRepositoryContract interface {
	FindAll() ([]domain.Barang, error)
	CheckBarang(nama_barang string) error
	InputBarang(data domain.Barang) (domain.Barang, *gorm.DB, error)
	InputBarangUpdate(data domain.Barang, tx *gorm.DB) error
	SearchBarang(column string, wildcardkey interface{}) ([]domain.Barang, error)
	DeleteBarang(barangID uint) (*gorm.DB, error)
	UpdateBarang(barang *domain.Barang, tx *gorm.DB) error
	BeginNewTransaction() *gorm.DB
}

type barangRepository struct {
	db *gorm.DB
}

func NewRepositoryBarang(db *gorm.DB) BarangRepositoryContract {
	return &barangRepository{db: db}
}  

func (br *barangRepository) BeginNewTransaction() *gorm.DB {
	return br.db.Begin()
}

func (br *barangRepository) FindAll() ([]domain.Barang, error) {
	var barang []domain.Barang

	if err := br.db.Find(&barang).Error; err != nil {
		return nil, err
	}

	return barang, nil
}

func (br *barangRepository) CheckBarang(nama_barang string) (error) {
	var barang domain.Barang

	if err := br.db.Model(&domain.Barang{}).Where("nama_barang = ?", nama_barang).First(&barang).Error; err != nil {
		return err
	}

	return nil
}

func (br *barangRepository) InputBarang(data domain.Barang) (domain.Barang, *gorm.DB, error) {
	tx := br.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		return domain.Barang{}, nil, err 
	}

	return data, tx, nil
}

// lanjutan dari proses inputBarang sehingga menggunakan tx dari InputBarang
func (br *barangRepository) InputBarangUpdate(data domain.Barang, tx *gorm.DB) error {
	if err := tx.Save(&data).Error; err != nil {
		return err
	}
	return nil
}

func (br *barangRepository) SearchBarang(column string, wildcardkey interface{}) ([]domain.Barang, error) {
	// mencari data berdasarkan key ke database
	var Barang []domain.Barang
	query := fmt.Sprintf("%s like ?", column)
	if err := br.db.Where(query, wildcardkey).Find(&Barang).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return []domain.Barang{}, errors.New("Barang Tidak Di Temukan")
		default:
			return []domain.Barang{}, err
		}
	}

	return Barang, nil
}


func (br *barangRepository) DeleteBarang(barangID uint) (*gorm.DB, error) {
	tx := br.db.Begin()

	if err := tx.Unscoped().Delete(&domain.Barang{}, barangID).Error; err != nil {
		return nil, err
	}

	return tx, nil
}

func (br *barangRepository) UpdateBarang(barang *domain.Barang, tx *gorm.DB) error {
	
	if err := tx.Model(&domain.Barang{}).Where("id = ?", barang.ID).Updates(&barang).Error; err != nil {
		return err
	}

	return nil
}

