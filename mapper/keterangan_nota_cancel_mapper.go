package mapper

import (
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
)

func ToKeteranganNotaCancelDTO(data domain.KeteranganNotaCancel) dto.KeteranganNotaCancel {
	return dto.KeteranganNotaCancel{
		ID: data.ID,
		Desc: data.Desc,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func ToKeteranganNotaCancelDomain(data *dto.KeteranganNotaCancel) domain.KeteranganNotaCancel {
	return domain.KeteranganNotaCancel{
		ID: data.ID,
		Desc: data.Desc,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}