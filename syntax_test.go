package mailvalidation_test

import (
	"testing"

	"github.com/pcasaretto/mailvalidation"
)

func TestSyntaxValidate(t *testing.T) {
	tests := []struct {
		In  string
		Out bool
	}{
		{"", false},
		{"pcasaretto@gmail.com", true},
		{"Paul <pcasaretto@gmail.com>", true},
	}

	validator := mailvalidation.Syntax{}
	for i, test := range tests {
		if actual, expected := validator.Validate(test.In), test.Out; actual != expected {
			t.Errorf("test %d\n\rgot: %t\n\rexpected: %t", i, actual, expected)
		}
	}
}
