package repository

import (
	"fmt"
	"net/http"
	"gorm.io/gorm"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
)

type HadiahRepository interface {
	BeginNewTransaction() *gorm.DB
	GetAllHadiah() ([]domain.Hadiah, error)
	CheckHadiah(data *domain.Hadiah) error
	InputHadiah(data *domain.Hadiah, tx *gorm.DB) error
	InputHadiahUpdate(data *domain.Hadiah, tx *gorm.DB) error
	DeleteHadiah(data *domain.Hadiah, tx *gorm.DB) error
	SearchHadiah(column string, wildcardkey interface{}, response map[string]interface{}) ([]domain.Hadiah, int, error)
	UpdateHadiah(data *domain.Hadiah, tx *gorm.DB) error
}

type hadiahRepository struct {
	db *gorm.DB
}

func NewRepositoryHadiah(db *gorm.DB) HadiahRepository {
	return &hadiahRepository{db: db}
}

func (hr *hadiahRepository) BeginNewTransaction() *gorm.DB {
	return hr.db.Begin()
}

func (hr *hadiahRepository) GetAllHadiah() ([]domain.Hadiah, error) {
	var gethadiah []domain.Hadiah

	if err := hr.db.Find(&gethadiah).Error; err != nil {
		return []domain.Hadiah{}, err
	}

	return gethadiah, nil
}

func (hr *hadiahRepository) CheckHadiah(data *domain.Hadiah) error {
	var exestingBarang domain.Hadiah

	if err := hr.db.Where("nama_barang = ?", data.Nama_Barang).First(&exestingBarang).Error; err != nil {
		return err
	}

	return nil
}

func (hr *hadiahRepository) InputHadiah(data *domain.Hadiah, tx *gorm.DB) error {
	// insert data ke database
	if err := tx.Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (hr *hadiahRepository) InputHadiahUpdate(data *domain.Hadiah, tx *gorm.DB) error {
	if err := tx.Save(data).Error; err != nil {
		return err
	}

	return nil
}

func (hr *hadiahRepository) DeleteHadiah(data *domain.Hadiah, tx *gorm.DB) error {
	if err := tx.Unscoped().Delete(&domain.Hadiah{}, data.ID).Error; err != nil {
		return err
	}

	return nil
}

func (hr *hadiahRepository) SearchHadiah(column string, wildcardkey interface{}, response map[string]interface{}) ([]domain.Hadiah, int, error) {
	var hadiah []domain.Hadiah

	query := fmt.Sprintf("%s like ?", column)
	if err := hr.db.Where(query, wildcardkey).Find(&hadiah).Error; err != nil {
		switch err{
		case gorm.ErrRecordNotFound:
			response["message"] = "Hadiah Tidak Di Temukan"
			return []domain.Hadiah{}, http.StatusBadRequest, err
		default:
			response["message"] = err.Error()
			return []domain.Hadiah{}, http.StatusInternalServerError, err
		}
	}

	return hadiah, http.StatusOK, nil
} 

func (hr *hadiahRepository) UpdateHadiah(data *domain.Hadiah, tx *gorm.DB) error {
	if err := tx.Model(&domain.Hadiah{}).Where("id = ?", data.ID).Updates(&data).Error; err != nil {
		return err
	}

	return nil
}