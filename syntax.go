package mailvalidation

import "net/mail"

type Syntax struct {
}

func (s Syntax) Validate(m string) bool {
	_, err := mail.ParseAddress(m)
	return err == nil
}
