package gmail

import "github.com/pcasaretto/mailvalidation/espgrammar"

var Validator *validator

const Domain = "gmail.com"

type validator struct {
}

func (v *validator) Validate(mail string) bool {
	return lex(mail)
}

func init() {
	Validator = &validator{}
	espgrammar.Register(Domain, Validator)
}
