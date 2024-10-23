package mapper

import (
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
)

// konversi from dto to domain
func ToPembelianDomain(data *dto.Pembelian) domain.Pembelian {
	return domain.Pembelian {	
		ID: data.ID,
		UserID: data.UserID,
		User: ToUserDomain(&data.User),
		CreatedAt: data.CreatedAt,
		Tanggal_Pembelian: data.Tanggal_Pembelian,
		Total_Harga: data.Total_Harga,
		Total_Keuntungan: data.Total_Keuntungan,
		Image: data.Image,
		KeteranganNotaCancelID: data.KeteranganNotaCancelID,
		Status: data.Status,
	}
}

func ToPembelianDTO(data *domain.Pembelian) dto.Pembelian {
	return dto.Pembelian {	
		ID: data.ID,
		UserID: data.UserID,
		User: ToUserDTO(data.User),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Tanggal_Pembelian: data.Tanggal_Pembelian,
		Total_Harga: data.Total_Harga,
		Total_Keuntungan: data.Total_Keuntungan,
		Image: data.Image,
		KeteranganNotaCancelID: data.KeteranganNotaCancelID,
		Status: data.Status,
		KeteranganNotaCancel: ToKeteranganNotaCancelDTO(data.KeteranganNotaCancel),
	}
}