package mailvalidation_test

import (
	"fmt"
	"net/mail"
	"testing"

	"github.com/pcasaretto/mailvalidation"
)

func TestDNSLookupValidatorValidate(t *testing.T) {
}

func TestValidate(t *testing.T) {
	validator := mailvalidation.NewDNSLookupValidator(nil)

	tests := []struct {
		in  string
		out bool
	}{
		{"dandrews0@delicious.com", true},
		{"lwashington1@geocities.jp", false},
		{"bthomas2@google.com.hk", true},
		{"lnelson3@liveinternet.ru", false},
		{"jhunter4@cnn.com", false},
		{"jhamilton5@adobe.com", true},
		{"calexander6@bloglines.com", false},
		{"eolson7@amazonaws.com", false},
		{"rmartinez8@businessinsider.com", false},
		{"jlynch9@arizona.edu", false},
	}

	for _, m := range tests {
		addr, _ := mail.ParseAddress(m.in)
		valid := validator.Validate(addr)

		if valid != m.out {
			fmt.Println(addr)
			t.Errorf("got %t, want %t", valid, m.out)
		}
	}
}
