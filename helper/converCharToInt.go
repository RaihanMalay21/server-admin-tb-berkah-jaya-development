package helper

import (
	"strings"
	"strconv"
)

func ConvertionToIntWithourChar(values string) (int, error) {
	// menghilang charcter titik di values 
	valueString := strings.ReplaceAll(values, ".", "")

	// men konversikan menjadi integer
	valueint, err := strconv.Atoi(valueString)
	if err != nil {
		return 0, err
	}

	return valueint, nil
}