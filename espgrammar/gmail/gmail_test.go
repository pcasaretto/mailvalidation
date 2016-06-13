package gmail_test

import "github.com/pcasaretto/mailvalidation/espgrammar/gmail"
import "testing"

func TestGmailValidate(t *testing.T) {
	tests := []struct {
		Email    string
		Expected bool
	}{
		{"robert", true},
		{"ROBERT", true},
		{"rob", false},
		{"robertroberroberroberroberroberroberroberroberroberrobertttttttttt", false},
		{".robert", false},
		{"robert.", false},
		{"robe$rt", false},
		{"robe..rt", false},
		{"robert+rob..bob", true},
		{"robert+robertroberroberroberroberroberroberroberroberroberrobertttttttttt", true},
	}
	for i, test := range tests {
		if actual := gmail.Validator.Validate(test.Email + "@" + gmail.Domain); actual != test.Expected {
			t.Errorf("test %d\n\texpected: %t\n\tgot: %t", i, test.Expected, actual)
		}
	}
}
