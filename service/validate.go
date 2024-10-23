package service

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/dto"
)

func ValidateStruct(data *dto.Barang) []map[string]string {
	validate := validator.New(validator.WithRequiredStructEnabled())
	Trans := helper.TranslatorIDN()

	var errors []map[string]string

	if err := validate.Struct(data); err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			NameField := err.StructField()
			errTranlate := err.Translate(Trans)
			errorMap := map[string]string{
				NameField: errTranlate,
			}
			errors = append(errors, errorMap)
		}

	}

	return errors
}

func ValidateStructHadiah(data *dto.Hadiah) []map[string]string {
	validate := validator.New(validator.WithRequiredStructEnabled())
	Trans := helper.TranslatorIDN()

	var errors []map[string]string

	if err := validate.Struct(data); err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			NameField := err.StructField()
			errTranlate := err.Translate(Trans)
			errorMap := map[string]string{
				NameField: errTranlate,
			}
			errors = append(errors, errorMap)
		}

	}

	return errors
}

func ValidateStructPembelianPerItem(data *dto.Pembelian_Per_Item, validate *validator.Validate, trans ut.Translator, response map[string]interface{}) error {
	if err := validate.Struct(data); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			response["message"] = "Invalid validation error: " + err.Error()
			return err
		}

		errors := err.(validator.ValidationErrors)
		errorsMessage := errors.Translate(trans)
		
		response["message"] = errorsMessage
		return err
	}

	return nil
} 

func ValidateStructPembelian(pembelian *dto.Pembelian, validate *validator.Validate, trans ut.Translator, response map[string]interface{}) error {
	if err := validate.Struct(pembelian); err != nil {
		// Check if the error is of type *validator.InvalidValidationError
		if _, ok := err.(*validator.InvalidValidationError); ok {
			// Handle InvalidValidationError
			response["message"] = "Invalid validation error:" + err.Error()
			return err
		}

		errors := err.(validator.ValidationErrors)
		errorsMessage := errors.Translate(trans)

		response["message"] = errorsMessage
		return err
	}

	return nil
}


func ValidateFileExtention(fileHeader *multipart.FileHeader, response map[string]interface{}) (string, error) {
	ext := filepath.Ext(fileHeader.Filename) // extention file image
	if err := helper.ValidateExtentionFile(ext); err != nil {
		response["message"] = err.Error()
		return "", err
	}
	return ext, nil
}

func ValidateFileSize(fileHeader *multipart.FileHeader, response map[string]interface{}) error {
	if fileHeader.Size > 2000000 {
		response["message"] = "Error file to large, max size file 2mb"
		return errors.New("Image to large")
	}
	return nil
}

func ValidateStuctBarang(barang *dto.Barang, response map[string]interface{}) error {
	errorStructs := ValidateStruct(barang)
	if len(errorStructs) > 0 {
		response["messageField"] = errorStructs
		return errors.New("Validation Errors")
	}
	return nil
}