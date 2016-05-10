package mailvalidation_test

import (
	"testing"

	"github.com/pcasaretto/mailvalidation"
)

type MockValidator struct {
	Result bool
}

func (v MockValidator) Validate(mail string) bool {
	return v.Result
}

func TestMultipleValidate(t *testing.T) {

	tests := []struct {
		Results  []bool
		Expected bool
	}{
		{[]bool{}, false},
		{[]bool{false}, false},
		{[]bool{true}, true},
		{[]bool{false, true}, false},
		{[]bool{false, false}, false},
		{[]bool{true, false}, false},
		{[]bool{true, true}, true},
	}

	for i, test := range tests {
		validator := mailvalidation.Multiple{}
		for _, b := range test.Results {
			validator = append(validator, MockValidator{b})
		}
		if actual, expected := validator.Validate(""), test.Expected; actual != expected {
			t.Errorf("test %d\n\rgot: %t\n\rexpected: %t", i, actual, expected)
		}
	}
}
