package mapper

import (
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
)

// konversi from domain to DTO 
func ToBarangDTO(barangDomain *domain.Barang) dto.Barang {
	return dto.Barang {
		ID: barangDomain.ID,
		CreatedAt: barangDomain.CreatedAt,
		UpdatedAt: barangDomain.UpdatedAt,
		Nama_Barang: barangDomain.Nama_Barang,
		Harga_Barang: barangDomain.Harga_Barang,
		Harga_Beli: barangDomain.Harga_Beli,
		Image: barangDomain.Image,
		Kode: barangDomain.Kode,
	}
}

// konversi from dto to domain
func ToBarangDomain(barangDto *dto.Barang) domain.Barang {
	return domain.Barang {
		ID: barangDto.ID,
		CreatedAt: barangDto.CreatedAt,
		UpdatedAt: barangDto.UpdatedAt,
		Nama_Barang: barangDto.Nama_Barang,
		Harga_Barang: barangDto.Harga_Barang,
		Harga_Beli: barangDto.Harga_Beli,
		Image: barangDto.Image,
		Kode: barangDto.Kode,
	}
}