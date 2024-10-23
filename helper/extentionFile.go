package helper

import (
	"errors"
)

func ValidateExtentionFile(ext string) error {
	if ext == "" || (ext != ".jpg" && ext != ".png" && ext != ".gift") {
		return errors.New("Invalid file extention: only .jpg, .png, and .gift only are allowed")
	}
	return nil
}