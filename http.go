package main

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

var client *http.Client

func init() {
	roundTriper := http.DefaultTransport
	transportPointer, ok := roundTriper.(*http.Transport)
	if !ok {
		panic("default roundtriper not an http.transport")
	}
	transport := *transportPointer
	transport.MaxIdleConns = 1000
	transport.MaxIdleConnsPerHost = 1000
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client = &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}
}

func getHTML(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return html, nil
}
