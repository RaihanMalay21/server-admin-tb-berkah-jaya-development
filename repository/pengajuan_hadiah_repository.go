package repository

import (
	"gorm.io/gorm"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
)

type PengajuanHadiahRepository interface {
	BeginNewTransaction() *gorm.DB 
	GetAllPengajuanHadiah() ([]domain.HadiahUser, error)
	GetDataUser(id uint) (domain.User, error)
	UpdateData(userID uint, hadiahID uint, GiftsArrive string, tx *gorm.DB) error
	UpdateDataHadiahFineshed(userID uint, hadiahID uint, status string) error
}

type pengajuanHadiahRepository struct {
	db *gorm.DB
}

func NewRepositoryPengajuanHadiah(db *gorm.DB) PengajuanHadiahRepository {
	return &pengajuanHadiahRepository{db: db}
}

func (repo *pengajuanHadiahRepository) BeginNewTransaction() *gorm.DB {
	return repo.db.Begin()
}

func (repo *pengajuanHadiahRepository) GetDataUser(id uint) (domain.User, error) {
	var user domain.User
	if err := repo.db.Where("id = ?", id).First(&user).Error; err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (repo *pengajuanHadiahRepository) GetAllPengajuanHadiah() ([]domain.HadiahUser, error) {
	var hadiahUser []domain.HadiahUser
	if err := repo.db.Preload("Hadiah").Preload("User").Where("gifts_arrive = ? or status = ?", "NO", "unfinished").Find(&hadiahUser).Error; err != nil {
		return nil, err
	}
	
	return hadiahUser, nil
}

func (repo *pengajuanHadiahRepository) UpdateData(userID uint, hadiahID uint, GiftsArrive string, tx *gorm.DB) error {
	if err := tx.Model(domain.HadiahUser{}).Where("user_id = ? and hadiah_id = ?", userID, hadiahID).Update("gifts_arrive", GiftsArrive).Error; err != nil {
		return err
	}

	return nil
}

func (repo *pengajuanHadiahRepository) UpdateDataHadiahFineshed(userID uint, hadiahID uint, status string) error {
	if err := repo.db.Model(domain.HadiahUser{}).Where("user_id = ? and hadiah_id = ?", userID, hadiahID).Update("status", status).Error; err != nil {
		return err
	}

	return nil
} 