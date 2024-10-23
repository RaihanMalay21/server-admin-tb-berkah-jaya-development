package repository

import (
	"gorm.io/gorm"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
)

type PoinRepository interface {
	BeginNewRepository() *gorm.DB
	SubmissionPoinGet() ([]domain.Pembelian, error)
	GetDataBarang(barang *domain.Barang) error
	UpdatePoinUser(amountPoin float64, userID uint, tx *gorm.DB) error
	UpdatePembelianPengajuanPoinSucces(pembelian *domain.Pembelian, tx *gorm.DB) error
	InputPembelianPerItem(pembelianperitem *domain.Pembelian_Per_Item, tx *gorm.DB) error 
	InputKeteranganNotaCancel(data *domain.KeteranganNotaCancel, tx *gorm.DB) error
	GetPoinUser(userID uint) (domain.User, error)
}

type poinRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) PoinRepository {
	return &poinRepository{db: db}
}

func (pr *poinRepository) BeginNewRepository() *gorm.DB {
	return pr.db.Begin()
}

func (pr *poinRepository) SubmissionPoinGet() ([]domain.Pembelian, error) {
	var pembelians []domain.Pembelian
	if err := pr.db.Where("status = 'cancel' and keterangan_nota_cancel_id IS NULL").Select("ID", "created_at", "updated_at", "user_id", "tanggal_pembelian", "total_harga", "image").Omit("keterangan_nota_cancel_id").Preload("User", func(db *gorm.DB) *gorm.DB{ return db.Select("id", "user_name", "email")}).Find(&pembelians).Error; err != nil {
		return nil, err
	}

	return pembelians, nil
}

func (pr *poinRepository) GetDataBarang(barang *domain.Barang) error {
	if err := pr.db.First(barang).Error; err != nil {
		return err
	}

	return nil
}

func (pr *poinRepository) GetPoinUser(userID uint) (domain.User, error) {
	var user domain.User
	if err := pr.db.Select("poin").Where("ID = ?", userID).First(&user).Error; err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (pr *poinRepository) UpdatePoinUser(amountPoin float64, userID uint, tx *gorm.DB) error {
	if err := tx.Model(&domain.User{}).Where("ID = ?", userID).Omit("keterangan_nota_cancel_id").Update("poin", amountPoin).Error; err != nil {
		return err
	}

	return nil
} 

func (pr *poinRepository) UpdatePembelianPengajuanPoinSucces(pembelian *domain.Pembelian, tx *gorm.DB) error {
	if err := tx.Model(&domain.Pembelian{}).Where("ID = ?", pembelian.ID).Updates(pembelian).Error; err != nil {
		return err
	}

	return nil
}

func (pr *poinRepository) InputPembelianPerItem(pembelianperitem *domain.Pembelian_Per_Item, tx *gorm.DB) error {
	if err := tx.Create(pembelianperitem).Error; err != nil {
		return err
	}

	return nil
}

func (pr *poinRepository) InputKeteranganNotaCancel(data *domain.KeteranganNotaCancel, tx *gorm.DB) error {
	if err := tx.Create(data).Error; err != nil {
		return err
	}

	return nil
}

