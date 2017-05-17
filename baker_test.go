package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

var bakerClient *http.Client

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func init() {
	transport := http.DefaultTransport.(*http.Transport)
	transport.MaxIdleConnsPerHost = 100

	bakerClient = &http.Client{
		Transport: transport,
	}
}

func BenchmarkGenQrcodeImg(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine has its own bytes.Buffer.
		for pb.Next() {
			// Example: this will give us a 44 byte, base64 encoded output
			token, err := GenerateRandomString(32)
			if err != nil {
				panic(err)
			} else {
				resp, err := bakerClient.Get("http://localhost:8080/merchant_qrcode?content=hehehe" + token)
				if err != nil {
					panic(err)
				} else {
					defer resp.Body.Close()
					io.Copy(ioutil.Discard, resp.Body)
				}
			}
		}
	})
}
