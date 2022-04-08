package restapi

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

func prepareSTSClientTransport(insecure bool) *http.Transport {
	DefaultTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 15 * time.Second,
		}).DialContext,
		MaxIdleConns:          1024,
		MaxIdleConnsPerHost:   1024,
		ResponseHeaderTimeout: 60 * time.Second,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    true,
		TLSClientConfig: &tls.Config{
			// Can't use SSLv3 because of POODLE and BEAST
			// Can't use TLSv1.0 because of POODLE and BEAST using CBC cipher
			// Can't use TLSv1.1 because of RC4 cipher usage
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: insecure,
			RootCAs:            GlobalRootCAs,
		},
	}
	return DefaultTransport
}

// PrepareConsoleHTTPClient returns an http.Client with custom configurations need it by *credentials.STSAssumeRole
// custom configurations include the use of CA certificates
func PrepareConsoleHTTPClient(insecure bool) *http.Client {
	transport := prepareSTSClientTransport(insecure)
	c := &http.Client{
		Transport: transport,
	}
	return c
}
