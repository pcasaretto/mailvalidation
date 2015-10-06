package main

import (
	"fmt"
	"net/mail"

	"github.com/k0kubun/pp"
	"github.com/pcasaretto/mailvalidation"
)

func __gore_p(xx ...interface{}) {
	for _, x := range xx {
		pp.Println(x)
	}
}
func main() {
	addr, _ := mail.ParseAddress("pcasaretto@gmail.com")
	validator := mailvalidation.NewDNSLookupValidator(nil)
	fmt.Println(validator.Validate(addr))
}
