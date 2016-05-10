package mailvalidation

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var LookupMX func(host string) ([]*net.MX, error) = net.LookupMX
var LookupIP func(host string) ([]net.IP, error) = net.LookupIP
var DialTimeout func(network, address string, timeout time.Duration) (net.Conn, error) = net.DialTimeout

type SMTPCheck struct {
	Timeout time.Duration
}

func NewSMTPCheck() *SMTPCheck {
	return &SMTPCheck{5 * time.Second}
}

func extractDomain(m string) string {
	return strings.SplitAfter(m, "@")[1]
}

type strategy func(domain string, hosts *[]string) bool

var checkMX strategy = func(domain string, hosts *[]string) bool {
	mxs, err := LookupMX(domain)
	if err != nil || len(mxs) == 0 {
		return false
	}
	for _, mx := range mxs {
		*hosts = append(*hosts, mx.Host)
	}
	return true
}

var checkA strategy = func(domain string, hosts *[]string) bool {
	ips, err := LookupIP(domain)
	if err != nil || len(ips) == 0 {
		return false
	} else {
		for _, ip := range ips {
			*hosts = append(*hosts, ip.String())
		}
	}
	return true
}

var strategies []strategy = []strategy{checkMX, checkA}

func (d *SMTPCheck) Validate(m string) bool {
	domain := extractDomain(m)

	hosts := make([]string, 0, 5)
	for _, s := range strategies {
		if s(domain, &hosts) {
			break
		}
	}
	if len(hosts) < 1 {
		return false
	}

	done := make(chan struct{})
	defer close(done)

	var outs []<-chan bool

	for _, host := range hosts {
		out := make(chan bool)
		outs = append(outs, out)
		go func(host string) {
			defer close(out)
			addr := fmt.Sprintf("%s:smtp", host)
			conn, err := DialTimeout("tcp", addr, d.Timeout)
			if err != nil || conn == nil {
				out <- false
				return
			}
			conn.Close()
			select {
			case out <- true:
			case <-done:
			}
		}(host)
	}

	total := len(outs)
	for {
		select {
		case r := <-merge(outs...):
			if r {
				return true
			}
			total--
			if total == 0 {
				return false
			}
		case <-time.After(d.Timeout):
			return false
		}
	}
}

func merge(cs ...<-chan bool) <-chan bool {
	var wg sync.WaitGroup
	out := make(chan bool)

	// Start an output goroutine for each input channel in cs.
	// output copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan bool) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
