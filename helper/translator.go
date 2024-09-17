package helper

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
)

func TranslatorIDN() ut.Translator {
	en := en.New()
	id := id.New()
	ut := ut.New(en, id)

	trans, _ := ut.GetTranslator("id")

	return trans
}