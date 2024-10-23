package mapper

import (
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
)

// konversi from dto to domain
func ToPengajuanHadiahDomain(data *dto.HadiahUser) domain.HadiahUser {
	return domain.HadiahUser {
		UserID: data.UserID,
		HadiahID: data.HadiahID,
		Hadiah: ToHadiahDomain(&data.Hadiah),
		User: ToUserDomain(&data.User),
		GiftsArrive: data.GiftsArrive,
		Status: data.Status,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

// konversi from domain to dto
func ToPengajuanHadiahDTO(data *domain.HadiahUser) dto.HadiahUser {
	return dto.HadiahUser {
		UserID: data.UserID,
		HadiahID: data.HadiahID,
		Hadiah: ToHadiahDTO(data.Hadiah),
		User: ToUserDTO(data.User),
		GiftsArrive: data.GiftsArrive,
		Status: data.Status,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}