package mapper

import (
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/domain"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
)

// konversi dari domain to dto
// field password di abaikan karna hal tersebut merupakan content sensitive 
// domain to dto yang akan di gunakan sebagai response request
func ToUserDTO(user domain.User) dto.User {
	return dto.User{
		ID: user.ID, 
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		UserName: user.UserName,
		Email: user.Email,
		NoWhatshapp: user.NoWhatshapp,
		Poin: user.Poin,
	}
}

// konversi from dto to domain
func ToUserDomain(user *dto.User) domain.User {
	return domain.User{
		ID: user.ID, 
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		UserName: user.UserName,
		Email: user.Email,
		Password: user.Password,
		NoWhatshapp: user.NoWhatshapp,
		Poin: user.Poin,
	}
}