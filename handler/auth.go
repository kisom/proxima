package handler

import (
	"crypto/subtle"
	"io"
	"log"
	"net/http"
	"time"

	"git.sr.ht/~kisom/goutils/config"
)

const (
	kHTTPUser     = "HTTP_USER"
	kHTTPPassword = "HTTP_PASSWORD"
)

// WithBasicAuth posts the data in the reader.
func WithBasicAuth(url string, contentType string, r io.Reader) (*http.Response, error) {
	transport := &http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.SetBasicAuth(config.Get(kHTTPUser), config.Get(kHTTPPassword))
	return transport.Do(req)
}

// CheckBasicAuth checks the authentication on an HTTP request.
func CheckBasicAuth(req *http.Request) bool {
	ruser, rpassword, ok := req.BasicAuth()
	if !ok {
		log.Println("request isn't using basic auth")
		return false
	}

	user := config.Get(kHTTPUser)
	password := config.Get(kHTTPPassword)

	subtleUser := subtle.ConstantTimeCompare([]byte(user), []byte(ruser))
	subtlePassword := subtle.ConstantTimeCompare(
		[]byte(password),
		[]byte(rpassword),
	)

	return subtleUser&subtlePassword == 1
}
