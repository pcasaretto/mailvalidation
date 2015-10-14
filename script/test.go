package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"sync"
	"sync/atomic"

	"github.com/pcasaretto/mailvalidation"
)

var total, invalid int32

func main() {

	var r io.Reader

	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		r = f
	} else {
		r = os.Stdin
	}

	ch := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go work(ch, &wg)
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ch <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	close(ch)
	wg.Wait()

	fmt.Println(total, invalid)
	// addr, _ := mail.ParseAddress("pcasaretto@gmail.com")
	// validator := mailvalidation.NewDNSLookupValidator(nil)
	// fmt.Println(validator.Validate(addr))
}

func work(in <-chan string, wg *sync.WaitGroup) {
	validator := mailvalidation.NewDNSLookupValidator(nil)
	for s := range in {
		atomic.AddInt32(&total, 1)

		m, err := mail.ParseAddress(s)
		if err != nil {
			atomic.AddInt32(&invalid, 1)
			continue
		}
		if !validator.Validate(m) {
			atomic.AddInt32(&invalid, 1)
		}
	}
	wg.Done()
}
