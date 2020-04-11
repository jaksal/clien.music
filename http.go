package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
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

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client = &http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
		Jar:       jar,
	}
}

func get(path string, headerList map[string]string, urlParamList map[string]string) ([]byte, error) {
	if urlParamList != nil {
		q := url.Values{}
		for k, v := range urlParamList {
			q.Add(k, v)
		}
		path += "?" + q.Encode()
	}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	if headerList != nil {
		for key, val := range headerList {
			req.Header.Add(key, val)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func post(path string, headerList map[string]string, bodyList map[string]string) ([]byte, error) {
	// fmt.Printf("param path=%s method=%s header=%+v url_param=%+v body=%+v\n", path, method, headerList, urlParamList, bodyList)
	var buf io.Reader
	if bodyList != nil {
		data := url.Values{}
		for key, val := range bodyList {
			data.Add(key, val)
		}
		buf = bytes.NewBufferString(data.Encode())
	}

	req, err := http.NewRequest("POST", path, buf)
	if err != nil {
		return nil, err
	}
	if headerList != nil {
		for key, val := range headerList {
			req.Header.Add(key, val)
		}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func login(uid, pwd string) error {
	p, err := get(origin+"/service/auth/login", nil, nil)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(p))
	if err != nil {
		return err
	}

	csrf, _ := doc.Find("input[name=_csrf]").Attr("value")
	log.Println(csrf)

	u := origin + "/service/login"
	if resp, err := post(u, map[string]string{
		//		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		//		"Accept-Encoding":           "deflate, br",
		//		"Accept-Language":           "ko,en;q=0.9",
		//		"Cache-Control":             "max-age=0",
		//		"Connection":                "keep-alive",
		//		"Host":                      "m.clien.net",
		//		"Origin":                    "https://m.clien.net",
		//		"Referer":                   "https://m.clien.net/service/auth/login",
		//		"Sec-Fetch-Dest":            "document",
		//		"Sec-Fetch-Mode":            "navigate",
		//		"Sec-Fetch-Site":            "same-origin",
		//		"Sec-Fetch-User":            "?1",
		//		"Upgrade-Insecure-Requests": "1",
		//		"User-Agent":                "Mozilla/5.0 (X11; CrOS x86_64 12871.67.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.94 Safari/537.36",
	}, map[string]string{
		"userId":       uid,
		"userPassword": pwd,
		"_csrf":        csrf,
	}); err != nil {
		panic(err)
	} else {

		fmt.Println(string(resp))

	}

	return nil
}
