package mailvalidation_test

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/pcasaretto/mailvalidation"
)

func TestValidate(t *testing.T) {

	goodConnection := func(network, address string, timeout time.Duration) (net.Conn, error) {
		conn, _ := net.Pipe()
		return conn, nil
	}
	badConnection := func(network, address string, timeout time.Duration) (net.Conn, error) {
		return nil, errors.New("error")
	}

	tests := []struct {
		Description string
		LookupMX    func(host string) ([]*net.MX, error)
		LookupIP    func(host string) ([]net.IP, error)
		DialTimeout func(network, address string, timeout time.Duration) (net.Conn, error)
		Expected    bool
	}{
		{
			"All nils",
			func(host string) (mxs []*net.MX, err error) {
				return nil, nil
			},
			func(host string) (ips []net.IP, err error) {
				return nil, nil
			},
			goodConnection,
			false,
		},
		{
			"MX ok",
			func(host string) (mxs []*net.MX, err error) {
				return []*net.MX{&net.MX{}}, nil
			},
			func(host string) (ips []net.IP, err error) {
				return nil, nil
			},
			goodConnection,
			true,
		},
		{
			"MX ok, bad connection",
			func(host string) (mxs []*net.MX, err error) {
				return []*net.MX{&net.MX{}, &net.MX{}}, nil
			},
			func(host string) (ips []net.IP, err error) {
				return nil, nil
			},
			badConnection,
			false,
		},
		{
			"MX error and IP ok",
			func(host string) (mxs []*net.MX, err error) {
				return nil, errors.New("error")
			},
			func(host string) (ips []net.IP, err error) {
				return []net.IP{net.IP{}, net.IP{}}, nil
			},
			goodConnection,
			true,
		},
		{
			"MX error and IP ok, bad connection",
			func(host string) (mxs []*net.MX, err error) {
				return nil, errors.New("error")
			},
			func(host string) (ips []net.IP, err error) {
				return []net.IP{net.IP{}}, nil
			},
			badConnection,
			false,
		},
	}

	for i, test := range tests {

		mailvalidation.LookupMX = test.LookupMX
		mailvalidation.LookupIP = test.LookupIP
		mailvalidation.DialTimeout = test.DialTimeout
		validator := mailvalidation.NewSMTPCheck()
		valid := validator.Validate("pcasaretto@gmail.com")

		if valid != test.Expected {
			t.Errorf("test %d:\n\rgot %t\n\rwant %t", i, valid, test.Expected)
		}
	}
}

func TestValidateTimeout(t *testing.T) {
	mailvalidation.LookupMX = func(host string) (mxs []*net.MX, err error) {
		return []*net.MX{&net.MX{}}, nil
	}
	mailvalidation.LookupIP = func(host string) (ips []net.IP, err error) {
		return nil, nil
	}
	mailvalidation.DialTimeout = func(network, address string, timeout time.Duration) (net.Conn, error) {
		time.Sleep(100 * time.Nanosecond)
		return nil, nil
	}

	validator := mailvalidation.NewSMTPCheck()
	validator.Timeout = time.Nanosecond

	valid := validator.Validate("pcasaretto@gmail.com")
	if valid {
		t.Errorf("got %t\n\rwant %t", valid, false)
	}
}
