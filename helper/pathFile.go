package helper

import "path/filepath"

func DestinationFolder(pathFolder, nameFile string) string {

	// menentukan path dir untuk file yang akan di create 
	filePath := filepath.Join(pathFolder, nameFile)

	return filePath
}