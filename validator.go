package mailvalidation

type Validator interface {
	Validate(mail string) bool
}
