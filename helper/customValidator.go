package helper

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"errors"

	config "github.com/RaihanMalay21/config-tb-berkah-jaya-development"
	models "github.com/RaihanMalay21/models_TB_Berkah_Jaya"
)

func isPhoneUnique(phone string) bool {
	// check apakah data exist in database 
	var user models.User
	if err := config.DB.Where("no_whatshapp = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		}
		return false
	}

	return false
}

func isEmailUnique(email string) bool {
	// check apakah data alredy exist in database
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		}
		return false
	}
	
	return false
}

func isUsernameUnique(username string) bool {
	// check apakah data alredy exist in database
	var user models.User
	if err := config.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		}
		return false
	}
	return false
}

func RegisterCustomValidations(validate *validator.Validate, trans ut.Translator) {
	validate.RegisterValidation("uniquePhone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		return isPhoneUnique(phone)
	})

	validate.RegisterTranslation("uniquePhone", trans, func(ut ut.Translator) error {
		return ut.Add("uniquePhone", "{0} sudah terdaftar", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("uniquePhone", fe.Field())
		return t
	})

	validate.RegisterValidation("uniqueEmail", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()
		return isEmailUnique(email)
	})

	validate.RegisterTranslation("uniqueEmail", trans, func(ut ut.Translator) error {
		return ut.Add("uniqueEmail", "{0} sudah terdaftar", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("uniqueEmail", fe.Field())
		return t
	})

	validate.RegisterValidation("uniqueUsername", func(fl validator.FieldLevel) bool {
		username := fl.Field().String()
		return isUsernameUnique(username)
	})

	validate.RegisterTranslation("uniqueUsername", trans, func(ut ut.Translator) error {
		return ut.Add("uniqueUsername", "{0} sudah terdaftar", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("uniqueUsername", fe.Field())
		return t
	})
}