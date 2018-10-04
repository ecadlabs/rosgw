package utils

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

	"gopkg.in/routeros.v2"
)

const (
	tcpScheme = "tcp"
	tlsScheme = "tls"
)

const (
	tcpPort = "8728"
	tlsPort = "8729"
)

type DialOptions struct {
	URL       string
	Username  string
	Password  string
	UseTLS    bool
	TLSConfig *tls.Config
}

func Dial(options *DialOptions) (*routeros.Client, error) {
	opt := options
	var addr string
	if u, err := url.Parse(opt.URL); err == nil && u.Host != "" {
		o := *opt
		opt = &o

		if u.Scheme == tlsScheme {
			opt.UseTLS = true
		} else if u.Scheme != tcpScheme {
			return nil, fmt.Errorf("Unknown URL scheme: %s", u.Scheme)
		}

		if ui := u.User; ui != nil {
			if un := ui.Username(); un != "" {
				opt.Username = un
			}

			if pw, ok := ui.Password(); ok {
				opt.Password = pw
			}
		}

		addr = u.Host
	} else {
		addr = opt.URL
	}

	if _, _, err := net.SplitHostPort(addr); err != nil {
		var port string
		if opt.UseTLS {
			port = tlsPort
		} else {
			port = tcpPort
		}

		addr = net.JoinHostPort(addr, port)
	}

	// TODO: Context support

	if opt.UseTLS {
		return routeros.DialTLS(addr, opt.Username, opt.Password, opt.TLSConfig)
	}

	return routeros.Dial(addr, opt.Username, opt.Password)
}
