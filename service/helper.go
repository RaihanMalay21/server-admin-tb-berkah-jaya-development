package service

import (
	"mime/multipart"
	"strings"
	"io"
	"os"
	"github.com/RaihanMalay21/server-admin-tb-berkah-jaya-development/helper"
)


func HelperSaveFile(image string, file multipart.File, response map[string]interface{}) error {
	pathFile := helper.DestinationFolder("C:\\Users\\acer\\Documents\\project app\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", image)

	outfile, err := os.Create(pathFile)
	if err != nil {
		response["message"] = err.Error()
		return err
	}
	defer outfile.Close()

	if _, err := io.Copy(outfile, file); err != nil {
		response["message"] =  err.Error()
		os.Remove(pathFile)
		return err
	}

	return nil
}

func DetetionWilcardAndColumn(key string, response map[string]interface{}) (string, interface{}) {
	// variabel untuk menampung field column wildcardkey
	var Column string
	var wildcardkey interface{}
	// memeriksa keynya 
	if strings.Contains(key, "."){
		Column = "harga_barang"
		// menghilangkan karakter titik and konversi ke int pda harga barang
		harga_barang, err := helper.ConvertionToIntWithourChar(key)
		if err != nil {
			response["message"] = err.Error()
			return "", nil
		}
		hargaBarangFloat64 := float64(harga_barang)
		wildcardkey = hargaBarangFloat64
		return Column, wildcardkey
	} else if helper.IsAllDigits(key){
		Column = "kode"
		wildcardkey = key + "%"
		return Column, wildcardkey
	} 
		
	Column = "nama_barang"
	wildcardkey = key + "%"

	return Column, wildcardkey
}

func RemoveFile(image string, response map[string]interface{}) error {
	filePath := helper.DestinationFolder("C:\\Users\\acer\\Documents\\project app\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", image)

	// melakukan penghapusan gambar
	if err := os.Remove(filePath); err != nil {
		response["message"] = err.Error()
		return err
	}

	return nil
} 

func RenameFile(image string, nameImageTemp string, response map[string]interface{}) error {
	// mengambil namanya saja untuk digabungkan dengan ext file yang di upload 
	beforeImagePath := helper.DestinationFolder("C:\\Users\\acer\\Documents\\project app\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", nameImageTemp)
	finalImagePath := helper.DestinationFolder("C:\\Users\\acer\\Documents\\project app\\development web berkah jaya\\fe_TB_Berkah_Jaya\\src\\images", image)

	// rename image tempory dengan nama image sebelumnya
	if err := os.Rename(beforeImagePath, finalImagePath); err != nil {
		response["message"] = err.Error()
		return err
	}

	return nil
} 





