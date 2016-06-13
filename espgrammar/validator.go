package espgrammar

import (
	"strings"

	"github.com/pcasaretto/mailvalidation"
)

var Validator *validator

type validator struct {
	validators map[string]mailvalidation.Validator
}

func (v *validator) Validate(mail string) bool {
	domain := extractDomain(mail)
	esp, ok := v.validators[domain]
	if !ok {
		return true
	}
	return esp.Validate(mail)
}

func init() {
	Validator = &validator{make(map[string]mailvalidation.Validator, 0)}
}

func Register(domain string, v mailvalidation.Validator) {
	Validator.validators[domain] = v
}

func extractDomain(m string) string {
	return strings.SplitAfter(m, "@")[1]
}
