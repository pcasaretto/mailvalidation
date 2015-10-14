package mailvalidation

import (
	"fmt"
	"net"
	"net/mail"
	"strings"
	"sync"
	"time"
)

type Validator interface {
	Validate(*mail.Address) (bool, error)
}

type DNSValidatorClient interface {
	LookupMX(host string) (mxs []*net.MX, err error)
	LookupIP(host string) (ips []net.IP, err error)
}

type defaultValidatorClient struct {
}

func (c defaultValidatorClient) LookupMX(host string) (mxs []*net.MX, err error) {
	return net.LookupMX(host)
}

func (c defaultValidatorClient) LookupIP(host string) (ips []net.IP, err error) {
	return net.LookupIP(host)
}

type DNSLookupValidator struct {
	dnsClient DNSValidatorClient
	Timeout   time.Duration
}

func NewDNSLookupValidator(client DNSValidatorClient) *DNSLookupValidator {
	if client == nil {
		client = defaultValidatorClient{}
	}
	return &DNSLookupValidator{client, 5 * time.Second}
}

func extractDomain(m *mail.Address) string {
	return strings.SplitAfter(m.Address, "@")[1]
}

func (d *DNSLookupValidator) Validate(m *mail.Address) bool {
	var hosts []string
	domain := extractDomain(m)

	// LookupMX
	mxs, err := d.dnsClient.LookupMX(domain)

	if err != nil || len(mxs) == 0 {
		// Lookup A
		ips, err := d.dnsClient.LookupIP(domain)
		fmt.Println("Ips: ", ips, err)
		if err != nil || len(ips) == 0 {
			return false
		} else {
			for _, ip := range ips {
				hosts = append(hosts, ip.String())
			}
		}
	} else {
		for _, mx := range mxs {
			hosts = append(hosts, mx.Host)
		}
	}
	// fmt.Println("Hosts: ", hosts)

	done := make(chan struct{})
	defer close(done)

	var outs []<-chan struct{}

	for _, host := range hosts {
		out := make(chan struct{})
		outs = append(outs, out)
		go func(host string) {
			addr := fmt.Sprintf("%s:smtp", host)
			// fmt.Println("dialing ", addr)
			conn, err := net.DialTimeout("tcp", addr, d.Timeout)
			if err != nil {
				return
			}
			conn.Close()
			select {
			case out <- struct{}{}:
			case <-done:
			}
		}(host)
	}

	select {
	case <-merge(outs...):
		return true
	case <-time.After(d.Timeout):
		return false
	}
}

func merge(cs ...<-chan struct{}) <-chan struct{} {
	var wg sync.WaitGroup
	out := make(chan struct{})

	// Start an output goroutine for each input channel in cs.
	// output copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan struct{}) {
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
