package mailvalidation

type Multiple []Validator

func (m Multiple) Validate(mail string) bool {
	if len(m) == 0 {
		return false
	}
	for _, v := range m {
		if !v.Validate(mail) {
			return false
		}
	}
	return true
}
